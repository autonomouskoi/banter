package banter

import (
	"context"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/twitch"
)

func (bb *Banterer) handleTwitchEvents(ctx context.Context) error {
	bb.bus.HandleTypes(ctx, twitch.BusTopics_TWITCH_EVENTSUB_EVENT.String(), 8,
		map[int32]bus.MessageHandler{
			int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_POINT_CUSTOM_REDEEM): bb.handleChannelPointRedeem,
		},
		nil)
	return nil
}

func (bb *Banterer) handleChannelPointRedeem(msg *bus.BusMessage) *bus.BusMessage {
	ecpcrr := &twitch.EventChannelPointsCustomRewardRedemption{}
	if err := bb.UnmarshalMessage(msg, ecpcrr); err != nil {
		return nil
	}
	bb.Log.Debug("channel point redeem", "ecpcrr", *ecpcrr)
	return nil
}
