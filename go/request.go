package banter

import (
	bus "github.com/autonomouskoi/core-tinygo"
)

// handle messages on the request topic
func (bb *Banterer) handleRequests() bus.TypeRouter {
	return bus.TypeRouter{
		int32(MessageTypeRequest_CONFIG_GET_REQ):    bb.handleRequestConfigGet,
		int32(MessageTypeRequest_BANTER_RENDER_REQ): bb.handleRequestBanterRender,
	}
}

// provide the config on request
func (bb *Banterer) handleRequestConfigGet(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	bus.MarshalMessage(reply, &ConfigGetResponse{
		Config: bb.cfg,
	})
	return reply
}

// provide a way for the user to test template rendering
func (bb *Banterer) handleRequestBanterRender(msg *bus.BusMessage) *bus.BusMessage {
	reply := bus.DefaultReply(msg)
	var req BanterRenderRequest
	if err := bus.UnmarshalMessage(msg, &req); err != nil {
		reply.Error = bus.InvalidTypeError(err)
		return reply
	}

	user := TestSender
	if req.Sender != nil {
		user = req.Sender
	}

	rendered, err := renderBanter(req.GetBanter(), req.GetOriginal(), user)
	if err != nil {
		reply.Error = bus.InvalidTypeError(err)
		return reply
	}

	bus.MarshalMessage(reply, &BanterRenderResponse{Output: rendered})
	return reply
}
