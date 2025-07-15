package banter

import (
	"strings"

	bus "github.com/autonomouskoi/core-tinygo"
	"github.com/autonomouskoi/twitch-tinygo"
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
		resp, err := twitch.GetUser(&twitch.GetUserRequest{Profile: bb.cfg.GetSendAs(), Login: part})
		if err != nil {
			bus.LogError("getting twitch user", "user", part, "error", err.Error())
			continue
		}
		msg := &bus.BusMessage{
			Topic: twitch.BusTopics_TWITCH_REQUEST.String(),
			Type:  int32(twitch.MessageTypeRequest_TYPE_REQUEST_SEND_SHOUTOUT_REQ),
		}
		bus.MarshalMessage(msg, &twitch.SendShoutoutRequest{
			FromProfile:     bb.cfg.GetSendAs(),
			ToBroadcasterId: resp.GetUser().Id,
			ModeratorId:     ccm.Chatter.Id,
		})
		if msg.Error != nil {
			return
		}
		bus.Send(msg)
	}
}
