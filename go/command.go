package banter

import (
	"sort"

	bus "github.com/autonomouskoi/core-tinygo"
)

// handle messages on the command topic
func (bb *Banterer) handleCommands() bus.TypeRouter {
	return bus.TypeRouter{
		int32(MessageTypeCommand_CONFIG_SET_REQ): bb.handleCommandConfigSet,
	}
}

// update the stored config
func (bb *Banterer) handleCommandConfigSet(msg *bus.BusMessage) *bus.BusMessage {
	reply := &bus.BusMessage{
		Topic: msg.GetTopic(),
		Type:  msg.Type + 1,
	}
	csr := &ConfigSetRequest{}
	if reply.Error = bus.UnmarshalMessage(msg, csr); reply.Error != nil {
		return reply
	}
	sort.Slice(csr.Config.Banters, func(i, j int) bool {
		return csr.Config.Banters[i].Command < csr.Config.Banters[j].Command
	})
	for _, banter := range csr.Config.Banters {
		if banter.Id == 0 {
			banter.Id = bb.rand.Int32()
		}
	}
	if err := csr.Config.Validate(); err != nil {
		reply.Error = &bus.Error{
			Code:   int32(bus.CommonErrorCode_INVALID_TYPE),
			Detail: bus.String(err.Error()),
		}
		return reply
	}
	bb.cfg = csr.GetConfig()
	bb.writeCfg()
	bus.MarshalMessage(reply, &ConfigSetResponse{
		Config: bb.cfg,
	})
	return reply
}
