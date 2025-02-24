package banter

import (
	"context"
	"strings"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/twitch"
)

func (bb *Banterer) sendIntegralCommand(cmd string, ccm *twitch.EventChannelChatMessage) {
	switch cmd {
	case "!so", "!shoutout":
		bb.sendShoutout(cmd, ccm)
	}
}

func (bb *Banterer) sendShoutout(cmd string, ccm *twitch.EventChannelChatMessage) {
	if !ccm.Enrichments.IsMod {
		return
	}
	args := strings.Split(strings.TrimPrefix(ccm.Message.Text, cmd), " ")
	for _, part := range args {
		if part == "" {
			continue
		}
		user, err := twitch.GetUser(context.Background(), bb.bus, bb.twitchProfile, part)
		if err != nil {
			bb.Log.Error("getting twitch user", "user", part, "error", err.Error())
			continue
		}
		msg := &bus.BusMessage{
			Topic: twitch.BusTopics_TWITCH_REQUEST.String(),
			Type:  int32(twitch.MessageTypeRequest_TYPE_REQUEST_SEND_SHOUTOUT_REQ),
		}
		bb.MarshalMessage(msg, &twitch.SendShoutoutRequest{
			FromProfile:     bb.twitchProfile,
			ToBroadcasterId: user.Id,
			ModeratorId:     ccm.Chatter.Id,
		})
		if msg.Error != nil {
			return
		}
		bb.bus.Send(msg)
	}
}
