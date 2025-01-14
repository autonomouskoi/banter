package banter

import (
	"bytes"
	"context"
	"text/template"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/twitch"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// handle events from Twitch, including chat messages, raids, follows, and cheers
func (bb *Banterer) handleTwitchEvents(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, twitch.BusTopics_TWITCH_EVENTSUB_EVENT.String(), 8,
		map[int32]bus.MessageHandler{
			int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_RAID):         bb.handleChannelRaid,
			int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_FOLLOW):       bb.handleChannelFollow,
			int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_CHEER):        bb.handleChannelCheer,
			int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_CHAT_MESSAGE): bb.handleChatMessageIn,
		},
		nil)
	return nil
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

// handle a raid, cheer, or follow using the appropriate template, if enabled
func (bb *Banterer) handleChannelEvent(msg *bus.BusMessage, settings *EventSettings, data protoreflect.ProtoMessage) {
	if !settings.Enabled || settings.Text == "" {
		return
	}
	tmpl, err := template.New("").Parse(settings.Text)
	if err != nil {
		bb.Log.Error("parsing template", "template", settings.Text, "error", err.Error())
		return
	}
	if err := bb.UnmarshalMessage(msg, data); err != nil {
		return
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		bb.Log.Error("executing template", "template", settings.Text, "error", err.Error())
		return
	}
	bb.sendChat(buf.String())
}
