// Package banter provides automation for Twitch chat.
package banter

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/autonomouskoi/akcore"
	bus "github.com/autonomouskoi/core-tinygo"
	"github.com/autonomouskoi/core-tinygo/svc"
	"github.com/autonomouskoi/datastruct/mapset"
	"github.com/autonomouskoi/twitch-tinygo"
)

var (
	topicTwitchRequest = twitch.BusTopics_TWITCH_REQUEST.String()
	cfgKVKey           = []byte("config")
)

// Banterer is the module
type Banterer struct {
	cfg       *Config
	cooldowns map[int32]int64 // random messages on cooldown aren't repeated
	router    bus.TopicRouter
	rand      *rand.Rand

	periodicNotifyToken int64
	lastPeriodicSend    int64
	guestListSeen       mapset.MapSet[string]
}

// Start banter
func New() (*Banterer, error) {
	bb := &Banterer{
		cooldowns:     map[int32]int64{},
		guestListSeen: mapset.MapSet[string]{},
	}

	var err error
	bb.periodicNotifyToken, err = svc.TimeNotifyEvery(1000)
	if err != nil {
		return nil, fmt.Errorf("requesting periodic notification: %w", err)
	}
	bb.rand = rand.New(rand.NewPCG(uint64(bb.periodicNotifyToken), uint64(time.Now().Unix())))

	if err := bb.loadConfig(); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	// we can't do anything without knowing which twitch profile we're going to
	// use. Wait for it to be ready
	bus.LogDebug("waiting for topic", "topic", topicTwitchRequest)
	for {
		hasTopic, err := bus.HasTopic(topicTwitchRequest, 1000)
		if err != nil {
			return nil, fmt.Errorf("waiting for topic %q: %w", topicTwitchRequest, err)
		}
		if hasTopic {
			break
		}
	}

	// pick a twitch profile
	msg := &bus.BusMessage{
		Topic: topicTwitchRequest,
		Type:  int32(twitch.MessageTypeRequest_TYPE_REQUEST_LIST_PROFILES_REQ),
	}
	if bus.MarshalMessage(msg, &twitch.ListProfilesRequest{}); msg.Error != nil {
		return nil, msg.Error
	}
	reply, err := bus.WaitForReply(msg, 1000)
	if err != nil {
		return nil, fmt.Errorf("waiting for twitch profiles: %w", err)
	}
	if reply.Error != nil {
		return nil, fmt.Errorf("listing twitch profiles: %v", reply.Error)
	}
	if bb.cfg.SendAs == "" {
		lpr := &twitch.ListProfilesResponse{}
		if err := bus.UnmarshalMessage(reply, lpr); err != nil {
			return nil, fmt.Errorf("unmarshalling: %w", err)
		}
		if len(lpr.GetNames()) == 0 {
			return nil, fmt.Errorf("no twitch profiles available")
		}
		for _, bb.cfg.SendAs = range lpr.GetNames() {
			break // just pick the first. In the future let the user pick
		}
		bb.writeCfg()
	}
	if bb.cfg.SendTo == "" {
		bb.cfg.SendTo = bb.cfg.SendAs
		bb.writeCfg()
	}

	bb.router = bus.TopicRouter{
		twitch.BusTopics_TWITCH_EVENTSUB_EVENT.String(): bb.handleTwitchEvents(),
		BusTopic_BANTER_REQUEST.String():                bb.handleRequests(),
		BusTopic_BANTER_COMMAND.String():                bb.handleCommands(),
		"9472ec79f0843765":                              bb.handleDirect(),
	}

	for topic := range bb.router {
		bus.LogDebug("subscribing", "topic", topic)
		if err := bus.Subscribe(topic); err != nil {
			return nil, fmt.Errorf("subscribing to topic %s: %w", topic, err)
		}
	}

	return bb, nil
}

func (bb *Banterer) Handle(msg *bus.BusMessage) {
	bb.router.Handle(msg)
}

func (bb *Banterer) loadConfig() error {
	bb.cfg = &Config{}
	if err := bus.KVGetProto(cfgKVKey, bb.cfg); err != nil && !errors.Is(err, akcore.ErrNotFound) {
		return fmt.Errorf("retrieving config: %w", err)
	}
	if bb.cfg.CooldownSeconds == 0 {
		bb.cfg.CooldownSeconds = 60 * 15
	}
	if bb.cfg.IntervalSeconds == 0 {
		bb.cfg.IntervalSeconds = 60 * 5
	}
	// this is only needed to migrate banters without ID numbers
	for _, banter := range bb.cfg.Banters {
		if banter.Id == 0 {
			banter.Id = bb.rand.Int32()
		}
	}
	return nil
}

func (bb *Banterer) writeCfg() {
	for _, banter := range bb.cfg.Banters {
		if banter.Id == 0 {
			banter.Id = bb.rand.Int32()
		}
	}
	if err := bus.KVSetProto(cfgKVKey, bb.cfg); err != nil {
		bus.LogError("writing config", "error", err.Error())
	}
}

func (bb *Banterer) handleDirect() bus.TypeRouter {
	return bus.TypeRouter{
		int32(svc.MessageType_TIME_NOTIFICATION_EVENT): bb.handleTimeNotificationEvent,
	}
}

// periodically send random messages
func (bb *Banterer) handleTimeNotificationEvent(msg *bus.BusMessage) *bus.BusMessage {
	var tn svc.TimeNotification
	if err := bus.UnmarshalMessage(msg, &tn); err != nil {
		bus.LogBusError("unmarshalling TimeNotification", err)
		return nil
	}
	if tn.GetToken() != bb.periodicNotifyToken {
		return nil
	}
	now := tn.GetCurrentTimeMillis()
	if bb.lastPeriodicSend+(1000*int64(bb.cfg.IntervalSeconds)) > now {
		return nil
	}
	bb.sendRandAnnouncement(now)
	bb.lastPeriodicSend = now
	return nil
}
