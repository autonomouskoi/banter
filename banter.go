package banter

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore"
	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules"
	"github.com/autonomouskoi/akcore/modules/modutil"
	"github.com/autonomouskoi/akcore/storage/kv"
	"github.com/autonomouskoi/akcore/web/webutil"
	"github.com/autonomouskoi/twitch"
)

const (
	EnvLocalContentPath = "AK_CONTENT_BANTER"
)

var (
	cfgKVKey = []byte("config")
)

func init() {
	manifest := &modules.Manifest{
		Id:          "9472ec79f0843765",
		Name:        "banter",
		Description: "Custom commands and periodic messages in Twitch Chat",
		WebPaths: []*modules.ManifestWebPath{
			{
				Path:        "/m/banter/",
				Type:        modules.ManifestWebPathType_MANIFEST_WEB_PATH_TYPE_EMBED_CONTROL,
				Description: "Controls for Banter",
			},
		},
	}
	modules.Register(manifest, &Banterer{})
}

//go:embed web.zip
var webZip []byte

type Banterer struct {
	http.Handler
	bus *bus.Bus
	modutil.ModuleBase
	lock      sync.Mutex
	kv        kv.KVPrefix
	cfg       *Config
	cooldowns map[string]time.Time
}

func (bb *Banterer) Start(ctx context.Context, deps *modutil.ModuleDeps) error {
	bb.bus = deps.Bus
	bb.Log = deps.Log
	bb.kv = deps.KV

	bb.cooldowns = map[string]time.Time{}

	if err := bb.loadConfig(); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	defer bb.writeCfg()

	fs, err := webutil.ZipOrEnvPath(EnvLocalContentPath, webZip)
	if err != nil {
		return fmt.Errorf("get web FS %w", err)
	}
	bb.Handler = http.StripPrefix("/m/banter", http.FileServer(fs))

	eg := errgroup.Group{}
	eg.Go(func() error { return bb.handleRequests(ctx) })
	eg.Go(func() error { return bb.handleCommands(ctx) })
	eg.Go(func() error { return bb.handleChat(ctx) })
	eg.Go(func() error { bb.periodicSend(ctx, time.Second); return nil })

	return eg.Wait()
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

func (bb *Banterer) handleRequests(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, BusTopic_BANTER_ANNOUNCE_REQUEST.String(), 8,
		map[int32]bus.MessageHandler{
			int32(MessageTypeRequest_CONFIG_GET_REQ): bb.handleRequestConfigGet,
		},
		nil,
	)
	return nil
}

func (bb *Banterer) handleRequestConfigGet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.Type + 1,
	}
	bb.lock.Lock()
	bb.MarshalMessage(reply, &ConfigGetResponse{
		Config: bb.cfg,
	})
	bb.lock.Unlock()
	return reply
}

func (bb *Banterer) handleCommands(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, BusTopic_BANTER_ANNOUNCE_COMMAND.String(), 4,
		map[int32]bus.MessageHandler{
			int32(MessageTypeCommand_CONFIG_SET_REQ): bb.handleCommandConfigSet,
		},
		nil,
	)
	return nil
}

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

func (bb *Banterer) handleChat(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, twitch.BusTopics_TWITCH_CHAT_EVENT.String(), 4,
		map[int32]bus.MessageHandler{
			int32(twitch.MessageTypeTwitchChatEvent_TWITCH_CHAT_EVENT_TYPE_MESSAGE_IN): bb.handleChatMessageIn,
		},
		nil,
	)
	return nil
}

func (bb *Banterer) handleChatMessageIn(msg *bus.BusMessage) *bus.BusMessage {
	cmi := &twitch.TwitchChatEventMessageIn{}
	if err := proto.Unmarshal(msg.GetMessage(), cmi); err != nil {
		bb.Log.Error("unmarshalling", "type", cmi.ProtoReflect().Descriptor().Name(), "error", err.Error())
		return nil
	}
	if !strings.HasPrefix(cmi.Text, "!") {
		return nil
	}
	text := strings.ToLower(cmi.Text)
	if text == "!banter" {
		bb.handleChatBanterList()
		return nil
	}
	var matchedBanter *Banter
	bb.lock.Lock()
	for _, banter := range bb.cfg.Banters {
		if banter.Command == text {
			matchedBanter = banter
			break
		}
	}
	bb.lock.Unlock()
	if matchedBanter != nil {
		bb.sendBanter(matchedBanter)
	}

	return nil
}

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

func (bb *Banterer) sendBanter(banter *Banter) {
	// TODO: banter text processing
	bb.sendChat(banter.Text)
	bb.cooldowns[banter.Command] = time.Now().Add(time.Second * time.Duration(bb.cfg.CooldownSeconds))
}

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

func (bb *Banterer) sendRandAnnouncement() {
	now := time.Now()
	var eligible []*Banter
	bb.lock.Lock()
	for _, banter := range bb.cfg.Banters {
		if banter.Disabled || !banter.Announce {
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
	bb.sendBanter(selected)
}

/*
func (s *spam) handle(msgIn *twitchchatpb.MessageIn) {
	text := msgIn.Text
	if !strings.HasPrefix(text, "!spam") {
		return
	}
	args := strings.Split(text, " ")[1:]
	if len(args) == 0 {
		twitchchat.Send(s.bus, "Available spams: %s", strings.Join(s.spamKeys, ", "))
		return
	}
	spamName := args[0]
	spam, present := s.spams[spamName]
	if !present {
		twitchchat.Send(s.bus, "Available spams: %s", strings.Join(s.spamKeys, ", "))
		return
	}
	s.cooldowns[spamName] = time.Now().Add(s.cooldown)
	twitchchat.Send(s.bus, "[bot] "+spam)
}
*/
