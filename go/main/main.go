package main

import (
	"github.com/extism/go-pdk"

	banter "github.com/autonomouskoi/banter/go"
	bus "github.com/autonomouskoi/core-tinygo"
)

var (
	b *banter.Banterer
)

//go:export start
func Initialize() int32 {
	bus.LogDebug("starting up")

	var err error
	b, err = banter.New()
	if err != nil {
		bus.LogError("creating banterer", "error", err.Error())
		return -1
	}

	bus.LogInfo("ready")

	return 0
}

//go:export recv
func Recv() int32 {
	msg := &bus.BusMessage{}
	if err := msg.UnmarshalVT(pdk.Input()); err != nil {
		bus.LogError("unmarshalling message", "error", err.Error())
		return 0
	}
	b.Handle(msg)
	return 0
}

func main() {}
