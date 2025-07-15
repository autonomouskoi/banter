package banter

import (
	protobuf_go_lite "github.com/aperturerobotics/protobuf-go-lite"

	bus "github.com/autonomouskoi/core-tinygo"
	"github.com/autonomouskoi/core-tinygo/svc"
	"github.com/autonomouskoi/twitch-tinygo"
)

// handle events from Twitch, including chat messages, raids, follows, and cheers
func (bb *Banterer) handleTwitchEvents() bus.TypeRouter {
	return bus.TypeRouter{
		int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_RAID):         bb.handleChannelRaid,
		int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_FOLLOW):       bb.handleChannelFollow,
		int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_CHEER):        bb.handleChannelCheer,
		int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_CHAT_MESSAGE): bb.handleChatMessageIn,
	}
}

func (bb *Banterer) handleChannelRaid(msg *bus.BusMessage) *bus.BusMessage {
	bb.handleChannelEvent(msg, bb.cfg.ChannelRaid, &twitch.EventChannelRaid{})
	return nil
}

func (bb *Banterer) handleChannelFollow(msg *bus.BusMessage) *bus.BusMessage {
	bb.handleChannelEvent(msg, bb.cfg.ChannelFollow, &twitch.EventChannelFollow{})
	return nil
}

func (bb *Banterer) handleChannelCheer(msg *bus.BusMessage) *bus.BusMessage {
	bb.handleChannelEvent(msg, bb.cfg.ChannelCheer, &twitch.EventChannelCheer{})
	return nil
}

type PB interface {
	protobuf_go_lite.Message
	protobuf_go_lite.JSONMessage
}

// handle a raid, cheer, or follow using the appropriate template, if enabled
func (bb *Banterer) handleChannelEvent(msg *bus.BusMessage, settings *EventSettings, data PB) {
	if !settings.Enabled || settings.Text == "" {
		return
	}
	json, err := data.MarshalJSON()
	if err != nil {
		bus.LogError("marshalling data to JSON", "err", err.Error())
		return
	}
	if err := bus.UnmarshalMessage(msg, data); err != nil {
		return
	}
	rendered, err := svc.RenderTemplate(settings.Text, json)
	if err != nil {
		bus.LogError("rendering template",
			"template", settings.Text, "json", string(json),
			"error", err.Error(),
		)
		return
	}
	bb.sendChat(rendered)
}
