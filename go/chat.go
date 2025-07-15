package banter

import (
	"strings"

	bus "github.com/autonomouskoi/core-tinygo"
	"github.com/autonomouskoi/core-tinygo/svc"
	"github.com/autonomouskoi/twitch-tinygo"
)

// send a message to Twitch chat using the twitch module
func (bb *Banterer) sendChat(text string) {
	msg := &bus.BusMessage{
		Topic: twitch.BusTopics_TWITCH_CHAT_REQUEST.String(),
		Type:  int32(twitch.MessageTypeTwitchChatRequest_TWITCH_CHAT_REQUEST_TYPE_SEND_REQ),
	}
	bus.MarshalMessage(msg, &twitch.TwitchChatRequestSendRequest{
		Text:    text,
		Channel: bb.cfg.SendTo,
		Profile: bb.cfg.SendAs,
	})
	reply, err := bus.WaitForReply(msg, 5000)
	if err != nil {
		bus.LogError("waiting for chat reply", "error", err.Error())
	}
	if reply.Error != nil {
		bus.LogError("waiting for chat reply", "error", reply.Error.Error())
	}
}

func renderBanter(banter *Banter, original *twitch.EventChannelChatMessage, sender *twitch.User) (string, error) {
	text := banter.GetText()
	bm, err := bus.ProtoMapJSON(map[string]bus.JSONMarshaller{
		"original":    original,
		"postCommand": bus.JSONString(strings.TrimSpace(strings.TrimPrefix(original.GetMessage().Text, banter.Command))),
		"sender":      sender,
	})

	if text, err = svc.RenderTemplate(text, bm); err != nil {
		bus.LogError("executing template", "command", banter.Command, "template", text, "error", err.Error())
		return "", err
	}
	return text, nil
}

// sendBanter composes and sends a banter chat message using the user-defined
// template and details about the message that triggered it
func (bb *Banterer) sendBanter(banter *Banter, original *twitch.EventChannelChatMessage) {
	text := banter.Text
	// if we have message details and it's a template it goes through processing
	if original != nil && strings.Contains(text, "{{") {
		bus.LogDebug("handling template message", "text", banter.Text)
		// get details about the chatter
		sender, err := twitch.GetUser(&twitch.GetUserRequest{
			Profile: bb.cfg.GetSendAs(),
			Login:   original.GetChatter().Name,
		})
		if err != nil {
			bus.LogError("getting twitch user",
				"login", original.GetChatter().Name,
				"error", err.Error(),
				"profile", bb.cfg.GetSendAs(),
			)
			return
		}
		if text, err = renderBanter(banter, original, sender.GetUser()); err != nil {
			bus.LogError("executing template", "command", banter.Command, "template", text, "error", err.Error())
			return
		}
		bus.LogDebug("processed message", "text", text)
	}
	bb.sendChat(text)
	now, _, err := svc.CurrentTimeMillis()
	if err != nil {
		bus.LogError("getting current time", "error", err.Error())
		return
	}
	// set a cooldown for this message
	bb.cooldowns[banter.Id] = now + (1000 * int64(bb.cfg.CooldownSeconds))
}

// handle a received chat message
func (bb *Banterer) handleChatMessageIn(msg *bus.BusMessage) *bus.BusMessage {
	ccm := &twitch.EventChannelChatMessage{}
	if err := bus.UnmarshalMessage(msg, ccm); err != nil {
		return nil
	}

	if !bb.guestListSeen.Has(ccm.Chatter.Id) {
		bb.guestListSeen.Add(ccm.Chatter.Id)
		bb.handleFirstSeen(ccm.Chatter)
	}

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
	for _, banter := range bb.cfg.Banters {
		if banter.Command == cmd {
			matchedBanter = banter
			break
		}
	}
	if matchedBanter != nil {
		bb.sendBanter(matchedBanter, ccm)
	} else {
		bb.sendIntegralCommand(cmd, ccm)
	}

	return nil
}

// handle the banter command to list banters
func (bb *Banterer) handleChatBanterList() {
	var commands []string
	for _, banter := range bb.cfg.Banters {
		if !banter.Disabled {
			commands = append(commands, banter.Command)
		}
	}
	bb.sendChat("Commands: " + strings.Join(commands, ", "))
}

// sendRandAnnouncements triggers a randomly-selected banter. Banters on
// cooldown or not marked Random are ineligible
func (bb *Banterer) sendRandAnnouncement(now int64) {
	var eligible []*Banter
	for _, banter := range bb.cfg.Banters {
		if banter.Disabled || !banter.Random {
			continue
		}
		if bb.cooldowns[banter.Id] > now {
			continue
		}
		eligible = append(eligible, banter)
	}

	if len(eligible) == 0 {
		return
	}

	idx := bb.rand.Int32N(int32(len(eligible)))
	selected := eligible[idx]
	bb.sendBanter(selected, nil)
}

func (bb *Banterer) handleFirstSeen(eventUser *twitch.EventUser) {
	var guestListCommands []*GuestListCommand
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
	if len(guestListCommands) == 0 {
		return
	}

	userResp, err := twitch.GetUser(&twitch.GetUserRequest{
		Profile: bb.cfg.GetSendAs(),
		Login:   eventUser.Login,
	})
	if err != nil {
		bus.LogError("getting twitch user",
			"for", "handleFirstSeen",
			"id", eventUser.Id,
			"login", eventUser.Login,
			"error", err.Error(),
		)
		return
	}

	for _, glc := range guestListCommands {
		userData, err := userResp.GetUser().MarshalJSON()
		if err != nil {
			bus.LogError("marshalling user to JSON", "error", err.Error())
			return
		}
		text, err := svc.RenderTemplate(glc.GetCommand(), userData)
		if err != nil {
			bus.LogError("rendering guest list command", "command", glc.Command, "error", err.Error())
			continue
		}
		if len(text) == 0 {
			return
		}
		bb.sendChat(text)
	}
}
