// Package banter provides automation for Twitch chat.
package banter

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	"github.com/autonomouskoi/akcore/web/webutil"
	"github.com/autonomouskoi/datastruct/mapset"
	"github.com/autonomouskoi/twitch"
)

const (
	// EnvLocalContentPath specifies an env var that, if set, provides a path
	// to serve web content from instead of the embedded zip
	EnvLocalContentPath = "AK_CONTENT_BANTER"
)

var (
	cfgKVKey           = []byte("config")
	topicTwitchRequest = twitch.BusTopics_TWITCH_REQUEST.String()
)

func init() {
	manifest := &modules.Manifest{
		Title:       "Banter",
		Id:          "9472ec79f0843765",
		Name:        "banter",
		Description: "Custom commands and periodic messages in Twitch Chat",
		WebPaths: []*modules.ManifestWebPath{
			{
				Path:        "https://autonomouskoi.org/module-banter.html",
				Type:        modules.ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_HELP,
				Description: "Help!",
			},
			{
				Path:        "/m/banter/embed_ctrl.js",
				Type:        modules.ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_EMBED_CONTROL,
				Description: "Controls for Banter",
			},
			{
				Path:        "/m/banter/index.html",
				Type:        modules.ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_CONTROL_PAGE,
				Description: "Controls for Banter",
			},
		},
	}
	modules.Register(manifest, &Banterer{})
}

//go:embed web.zip
var webZip []byte

// Banterer is the module
type Banterer struct {
	http.Handler
	bus *bus.Bus
	modutil.ModuleBase
	lock          sync.Mutex
	kv            kv.KVPrefix
	cfg           *Config
	cooldowns     map[string]time.Time // random messages on cooldown aren't repeated
	twitchProfile string

	guestListSeen mapset.MapSet[string]
}

// Start banter
func (bb *Banterer) Start(ctx context.Context, deps *modutil.ModuleDeps) error {
	bb.bus = deps.Bus
	bb.Log = deps.Log
	bb.kv = deps.KV

	bb.cooldowns = map[string]time.Time{}
	if bb.guestListSeen == nil {
		bb.guestListSeen = mapset.MapSet[string]{}
	}

	if err := bb.loadConfig(); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	defer bb.writeCfg()

	fs, err := webutil.ZipOrEnvPath(EnvLocalContentPath, webZip)
	if err != nil {
		return fmt.Errorf("get web FS %w", err)
	}
	bb.Handler = http.FileServer(fs)

	// we can't do anything without knowing which twitch profile we're going to
	// use. Wait for it to be ready
	bb.Log.Debug("waiting for topic", "topic", topicTwitchRequest)
	if err := bb.bus.WaitForTopic(ctx, topicTwitchRequest, time.Millisecond*10); err != nil {
		return fmt.Errorf("waiting for %s: %w", topicTwitchRequest, err)
	}

	// pick a twitch profile
	msg := &bus.BusMessage{
		Topic: topicTwitchRequest,
		Type:  int32(twitch.MessageTypeRequest_TYPE_REQUEST_LIST_PROFILES_REQ),
	}
	if bb.MarshalMessage(msg, &twitch.ListProfilesRequest{}); msg.Error != nil {
		return msg.Error
	}
	reply := bb.bus.WaitForReply(ctx, msg)
	if reply.Error != nil {
		return fmt.Errorf("listing twitch profiles: %w", err)
	}
	lpr := &twitch.ListProfilesResponse{}
	if err := bb.UnmarshalMessage(reply, lpr); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	if len(lpr.GetNames()) == 0 {
		return fmt.Errorf("no twitch profiles available")
	}
	for _, bb.twitchProfile = range lpr.GetNames() {
		break // just pick the first. In the future let the user pick
	}

	bb.Go(func() error { return bb.handleRequests(ctx) })
	bb.Go(func() error { return bb.handleCommands(ctx) })
	bb.Go(func() error { return bb.handleTwitchEvents(ctx) })
	bb.Go(func() error { bb.periodicSend(ctx, time.Second); return nil })

	return bb.Wait()
}

func (bb *Banterer) loadConfig() error {
	bb.cfg = &Config{}
	if err := bb.kv.GetProto(cfgKVKey, bb.cfg); err != nil && !errors.Is(err, akcore.ErrNotFound) {
		return fmt.Errorf("retrieving config: %w", err)
	}
	if bb.cfg.CooldownSeconds == 0 {
		bb.cfg.CooldownSeconds = 60 * 15
	}
	if bb.cfg.IntervalSeconds == 0 {
		bb.cfg.IntervalSeconds = 60 * 5
	}
	return nil
}

func (bb *Banterer) writeCfg() {
	bb.lock.Lock()
	defer bb.lock.Unlock()
	if err := bb.kv.SetProto(cfgKVKey, bb.cfg); err != nil {
		bb.Log.Error("writing config", "error", err.Error())
	}
}

// handle messages on the request topic
func (bb *Banterer) handleRequests(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, BusTopic_BANTER_REQUEST.String(), 8,
		map[int32]bus.MessageHandler{
			int32(MessageTypeRequest_CONFIG_GET_REQ): bb.handleRequestConfigGet,
		},
		nil,
	)
	return nil
}

// provide the config on request
func (bb *Banterer) handleRequestConfigGet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	bb.lock.Lock()
	bb.MarshalMessage(reply, &ConfigGetResponse{
		Config: bb.cfg,
	})
	bb.lock.Unlock()
	return reply
}

// handle messages on the command topic
func (bb *Banterer) handleCommands(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, BusTopic_BANTER_COMMAND.String(), 4,
		map[int32]bus.MessageHandler{
			int32(MessageTypeCommand_CONFIG_SET_REQ): bb.handleCommandConfigSet,
		},
		nil,
	)
	return nil
}

// update the stored config
func (bb *Banterer) handleCommandConfigSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.Type + 1,
	}
	csr := &ConfigSetRequest{}
	if reply.Error = bb.UnmarshalMessage(msg, csr); reply.Error != nil {
		return reply
	}
	sort.Slice(csr.Config.Banters, func(i, j int) bool {
		return csr.Config.Banters[i].Command < csr.Config.Banters[j].Command
	})
	if err := csr.Config.Validate(); err != nil {
		reply.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: proto.String(err.Error()),
		}
		return reply
	}
	bb.lock.Lock()
	bb.cfg = csr.GetConfig()
	bb.lock.Unlock()
	bb.writeCfg()
	bb.MarshalMessage(reply, &ConfigSetResponse{
		Config: bb.cfg,
	})
	return reply
}

// handle a received chat message
func (bb *Banterer) handleChatMessageIn(msg *bus.BusMessage) *bus.BusMessage {
	ccm := &twitch.EventChannelChatMessage{}
	if err := bb.UnmarshalMessage(msg, ccm); err != nil {
		return nil
	}

	bb.lock.Lock()
	if !bb.guestListSeen.Has(ccm.Chatter.Id) {
		go bb.handleFirstSeen(ccm.Chatter)
		bb.guestListSeen.Add(ccm.Chatter.Id)
	}
	bb.lock.Unlock()

	if !strings.HasPrefix(ccm.Message.Text, "!") {
		return nil
	}
	text := strings.ToLower(ccm.Message.Text)
	if text == "!banter" {
		bb.handleChatBanterList()
		return nil
	}
	cmd, _, _ := strings.Cut(text, " ")
	var matchedBanter *Banter
	bb.lock.Lock()
	for _, banter := range bb.cfg.Banters {
		if banter.Command == cmd {
			matchedBanter = banter
			break
		}
	}
	bb.lock.Unlock()
	if matchedBanter != nil {
		bb.sendBanter(matchedBanter, ccm)
	}

	return nil
}

// handle the banter command to list banters
func (bb *Banterer) handleChatBanterList() {
	var commands []string
	bb.lock.Lock()
	for _, banter := range bb.cfg.Banters {
		if !banter.Disabled {
			commands = append(commands, banter.Command)
		}
	}
	bb.lock.Unlock()
	bb.sendChat("Commands: " + strings.Join(commands, ", "))
}

// periodically send random messages
func (bb *Banterer) periodicSend(ctx context.Context, interval time.Duration) {
	lastSend := time.Now()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			if lastSend.Add(time.Second * time.Duration(bb.cfg.IntervalSeconds)).After(now) {
				continue
			}
			bb.sendRandAnnouncement()
			lastSend = now
		}
	}
}

// banterMessage carries details of a message to send to chat to be processed
// by the user's template
type banterMessage struct {
	Sender      *twitch.User
	Original    *twitch.EventChannelChatMessage
	PostCommand string
}

// sendBanter composes and sends a banter chat message using the user-defined
// template and details about the message that triggered it
func (bb *Banterer) sendBanter(banter *Banter, original *twitch.EventChannelChatMessage) {
	text := banter.Text
	// if we have message details and it's a template it goes through processing
	if original != nil && strings.Contains(text, "{{") {
		bb.Log.Debug("handling template message", "text", banter.Text)
		tmpl, err := template.New("").Parse(text)
		if err != nil {
			bb.Log.Error("parsing template", "command", banter.Command, "template", text, "error", err.Error())
			return
		}
		// get details about the chatter
		sender, err := twitch.GetUser(context.Background(), bb.bus, bb.twitchProfile, original.GetChatter().Name)
		if err != nil {
			bb.Log.Error("getting twitch user", "login", original.GetChatter().Name, "error", err.Error(), "profile", bb.twitchProfile)
			return
		}
		bm := banterMessage{
			Original:    original,
			PostCommand: strings.TrimSpace(strings.TrimPrefix(original.GetMessage().Text, banter.Command)),
			Sender:      sender,
		}
		buf := &bytes.Buffer{}
		if err := tmpl.Execute(buf, bm); err != nil {
			bb.Log.Error("executing template", "command", banter.Command, "template", text, "error", err.Error())
			return
		}
		text = buf.String()
		bb.Log.Debug("processed message", "text", text)
	}
	bb.sendChat(text)
	// set a cooldown for this message
	bb.cooldowns[banter.Command] = time.Now().Add(time.Second * time.Duration(bb.cfg.CooldownSeconds))
}

// send a message to Twitch chat using the twitch module
func (bb *Banterer) sendChat(text string) {
	msg := &bus.BusMessage{
		Topic: twitch.BusTopics_TWITCH_CHAT_REQUEST.String(),
		Type:  int32(twitch.MessageTypeTwitchChatRequest_TWITCH_CHAT_REQUEST_TYPE_SEND_REQ),
	}
	bb.MarshalMessage(msg, &twitch.TwitchChatRequestSendRequest{
		Text: text,
	})
	bb.bus.Send(msg)
}

// sendRandAnnouncements triggers a randomly-selected banter. Banters on
// cooldown or not marked Random are ineligible
func (bb *Banterer) sendRandAnnouncement() {
	now := time.Now()
	var eligible []*Banter
	bb.lock.Lock()
	for _, banter := range bb.cfg.Banters {
		if banter.Disabled || !banter.Random {
			continue
		}
		if bb.cooldowns[banter.Command].After(now) {
			continue
		}
		eligible = append(eligible, banter)
	}
	bb.lock.Unlock()

	if len(eligible) == 0 {
		return
	}

	idx := rand.Int31n(int32(len(eligible)))
	selected := eligible[idx]
	bb.sendBanter(selected, nil)
}

func (bb *Banterer) handleFirstSeen(eventUser *twitch.EventUser) {
	var guestListCommands []*GuestListCommand
	bb.lock.Lock()
COMMAND:
	for _, glc := range bb.cfg.GuestListCommands {
		for _, listName := range glc.GuestListNames {
			gl := bb.cfg.GuestLists[listName]
			if gl == nil {
				continue
			}
			for _, listMember := range gl.Members {
				if listMember.Id == eventUser.Id {
					guestListCommands = append(guestListCommands, glc)
					continue COMMAND
				}
			}
		}
	}
	bb.lock.Unlock()
	if len(guestListCommands) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	user, err := twitch.GetUser(ctx, bb.bus, bb.twitchProfile, eventUser.Login)
	if err != nil {
		bb.Log.Error("getting twitch user",
			"id", eventUser.Id,
			"login", eventUser.Login,
			"error", err.Error(),
		)
		return
	}

	buf := &bytes.Buffer{}
	for _, glc := range guestListCommands {
		buf.Reset()
		if err := renderGuestListCommand(glc, user, buf); err != nil {
			bb.Log.Error("rendering guest list command", "command", glc.Command, "error", err.Error())
			continue
		}
		text := strings.TrimSpace(buf.String())
		if len(text) == 0 {
			return
		}
		bb.sendChat(buf.String())
	}
}

func renderGuestListCommand(glc *GuestListCommand, user *twitch.User, w io.Writer) error {
	tmpl, err := template.New("").Parse(glc.GetCommand())
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}
	return tmpl.Execute(w, user)
}

//go:embed icon.svg
var icon []byte

func (*Banterer) Icon() ([]byte, string, error) {
	return icon, "image/svg+xml", nil
}
