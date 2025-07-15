package banter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/bus/bustest"
	banter "github.com/autonomouskoi/banter/go"
	"github.com/autonomouskoi/twitch"
)

var (
	topicBanterCommand     = banter.BusTopic_BANTER_COMMAND.String()
	topicBanterRequest     = banter.BusTopic_BANTER_REQUEST.String()
	topicTwitchRequest     = twitch.BusTopics_TWITCH_REQUEST.String()
	topicTwitchEventsub    = twitch.BusTopics_TWITCH_EVENTSUB_EVENT.String()
	topicTwitchChatRequest = twitch.BusTopics_TWITCH_CHAT_REQUEST.String()
)

func testListTestProfile(msg *bus.BusMessage) *bus.BusMessage {
	lpr := &twitch.ListProfilesResponse{Names: []string{"test-profile"}}
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.GetType() + 1,
	}
	reply.Message, _ = proto.Marshal(lpr)
	return reply
}

func TestConfig(t *testing.T) {
	t.Parallel()
	ctx, cancel, deps := bustest.NewDeps(t)

	eg := errgroup.Group{}

	eg.Go(func() error {
		deps.Bus.HandleTypes(ctx, topicTwitchRequest, 1,
			map[int32]bus.MessageHandler{
				int32(twitch.MessageTypeRequest_TYPE_REQUEST_LIST_PROFILES_REQ): testListTestProfile,
			}, nil)
		return nil
	})

	bb := &banter.Banterer{}
	eg.Go(func() error {
		return bb.Start(ctx, deps)
	})

	require.NoError(t, deps.Bus.WaitForTopic(ctx, topicBanterRequest, time.Millisecond), "waiting for banter")

	getCfgMsg := &bus.BusMessage{
		Topic: topicBanterRequest,
		Type:  int32(banter.MessageTypeRequest_CONFIG_GET_REQ),
	}
	getCfgMsg.Message, _ = proto.Marshal(&banter.ConfigGetRequest{})
	reply := deps.Bus.WaitForReply(ctx, getCfgMsg)
	require.Nil(t, reply.Error)
	cgResp := &banter.ConfigGetResponse{}
	require.NoError(t, proto.Unmarshal(reply.GetMessage(), cgResp), "unmarshalling cgResp")
	require.Equal(t, uint32(900), cgResp.Config.CooldownSeconds) // 900 is default

	// let's update it
	cgResp.Config.CooldownSeconds = 123
	setCfgMsg := &bus.BusMessage{
		Topic: topicBanterCommand,
		Type:  int32(banter.MessageTypeCommand_CONFIG_SET_REQ),
	}
	setCfgMsg.Message, _ = proto.Marshal(&banter.ConfigSetRequest{
		Config: cgResp.Config,
	})
	reply = deps.Bus.WaitForReply(ctx, setCfgMsg)
	require.Nil(t, reply.Error)
	csResp := &banter.ConfigSetResponse{}
	require.NoError(t, proto.Unmarshal(reply.GetMessage(), csResp), "unmarshalling csResp")
	require.Equal(t, uint32(123), csResp.Config.CooldownSeconds)

	reply = deps.Bus.WaitForReply(ctx, getCfgMsg)
	require.Nil(t, reply.Error)
	cgResp.Reset()
	require.NoError(t, proto.Unmarshal(reply.GetMessage(), cgResp), "unmarshalling cgResp")
	require.Equal(t, uint32(123), cgResp.Config.CooldownSeconds) // 900 is default

	cancel()
	require.NoError(t, eg.Wait())
}

func TestEvents(t *testing.T) {
	t.Parallel()
	ctx, cancel, deps := bustest.NewDeps(t)

	eg := errgroup.Group{}

	eg.Go(func() error {
		deps.Bus.HandleTypes(ctx, topicTwitchRequest, 1,
			map[int32]bus.MessageHandler{
				int32(twitch.MessageTypeRequest_TYPE_REQUEST_LIST_PROFILES_REQ): testListTestProfile,
			}, nil)
		return nil
	})

	bb := &banter.Banterer{}
	eg.Go(func() error {
		return bb.Start(ctx, deps)
	})

	require.NoError(t, deps.Bus.WaitForTopic(ctx, topicBanterCommand, time.Millisecond), "waiting for banter")
	cfg := &banter.Config{
		IntervalSeconds: 30,
		ChannelCheer: &banter.EventSettings{
			Enabled: true,
			Text:    "{{ .From.Name }} just cheered {{ .Bits }} bits!",
		},
		ChannelRaid: &banter.EventSettings{
			Enabled: true,
			Text:    "raid happened",
		},
	}
	msg := &bus.BusMessage{
		Topic: topicBanterCommand,
		Type:  int32(banter.MessageTypeCommand_CONFIG_SET_REQ),
	}
	msg.Message, _ = proto.Marshal(&banter.ConfigSetRequest{Config: cfg})
	reply := deps.Bus.WaitForReply(ctx, msg)
	require.Nil(t, reply.Error, reply.Error.GetDetail())

	require.NoError(t, deps.Bus.WaitForTopic(ctx, topicTwitchEventsub, time.Millisecond), "waiting for twitch")

	chat := make(chan *bus.BusMessage, 1)
	deps.Bus.Subscribe(topicTwitchChatRequest, chat)

	raidMsg := &bus.BusMessage{
		Topic: topicTwitchEventsub,
		Type:  int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_RAID),
	}
	var err error
	raidMsg.Message, err = proto.Marshal(&twitch.EventChannelRaid{})
	require.NoError(t, err, "marshalling raid")
	deps.Bus.Send(raidMsg)

	var chatMsg *bus.BusMessage
	select {
	case <-ctx.Done():
	case chatMsg = <-chat:
	}

	require.NotNil(t, chatMsg)
	crsr := &twitch.TwitchChatRequestSendRequest{}
	require.NoError(t, proto.Unmarshal(chatMsg.GetMessage(), crsr))
	require.Equal(t, "raid happened", crsr.GetText())

	cheerMsg := &bus.BusMessage{
		Topic: topicTwitchEventsub,
		Type:  int32(twitch.MessageTypeEventSub_TYPE_CHANNEL_CHEER),
	}
	cheerMsg.Message, err = proto.Marshal(&twitch.EventChannelCheer{
		From: &twitch.EventUser{
			Name: "TheSender",
		},
		Bits: 20,
	})
	require.NoError(t, err, "marshalling cheer")
	deps.Bus.Send(cheerMsg)

	select {
	case <-ctx.Done():
	case chatMsg = <-chat:
	}
	require.NotNil(t, chatMsg)
	crsr = &twitch.TwitchChatRequestSendRequest{}
	require.NoError(t, proto.Unmarshal(chatMsg.GetMessage(), crsr))
	require.Equal(t, "TheSender just cheered 20 bits!", crsr.GetText())

	cancel()
	require.NoError(t, eg.Wait())
}
