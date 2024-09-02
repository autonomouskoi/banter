package banter_test

/*

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/autonomouskoi/akcore/bus"
	"github.com/autonomouskoi/akcore/modules/modutil"
)

func TestSpam(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	requireT := require.New(t)

	cfg := &config.Spam{
		Interval: time.Second / 4,
		Cooldown: time.Second,
		Spams: map[string]string{
			"blarg": "honk",
		},
	}

	deps := &modutil.Deps{
		Bus: bus.New(ctx),
		Log: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	}

	chatOut := make(chan *akpb.BusMessage, 5)
	deps.Bus.Subscribe(twitchchatpb.BusTopics_TWITCHCHAT_SEND.String(), chatOut)
	err := deps.Bus.WaitForTopic(ctx, twitchchatpb.BusTopics_TWITCHCHAT_SEND.String(), time.Millisecond*10)
	requireT.NoError(err, "waiting for chat sender")

	go func() {
		requireT.NoError(spam.Start(cfg, deps), "starting spam")
	}()

	err = deps.Bus.WaitForTopic(ctx, twitchchatpb.BusTopics_TWITCHCHAT_RECV.String(), time.Millisecond*10)
	requireT.NoError(err, "waiting for chat receiver")

	sendChatMsg := func(s string) {
		t.Helper()
		m, err := proto.Marshal(&twitchchatpb.MessageIn{
			Text: s,
		})
		requireT.NoError(err, "marshalling chat message")
		deps.Bus.Send(
			&akpb.BusMessage{
				Topic:   twitchchatpb.BusTopics_TWITCHCHAT_RECV.String(),
				Message: m,
			})
	}

	parseMsgOut := func(msg *akpb.BusMessage) *twitchchatpb.MessageOut {
		t.Helper()
		requireT.NotNil(msg)
		mo := &twitchchatpb.MessageOut{}
		err := proto.Unmarshal(msg.GetMessage(), mo)
		requireT.NoError(err)
		return mo
	}

	// random message
	got := <-chatOut
	cm := parseMsgOut(got)
	requireT.Equal("[bot] honk", cm.Text)

	// basically "help"
	sendChatMsg("!spam")
	got = <-chatOut
	cm = parseMsgOut(got)
	requireT.Equal("Available spams: blarg", cm.Text)

	// basically "help"
	sendChatMsg("!spam bad-command")
	got = <-chatOut
	cm = parseMsgOut(got)
	requireT.Equal("Available spams: blarg", cm.Text)

	// This will be ignored
	sendChatMsg("spam test")

	// This will be handled
	sendChatMsg("!spam blarg")

	// the message we expect to handle
	got = <-chatOut
	cm = parseMsgOut(got)
	requireT.Equal("[bot] honk", cm.Text)

	// close the sources
	cancel()

	// no other messages
	for msg := range chatOut {
		t.Fatalf("unexpected message: %#v", msg)
	}
}
*/
