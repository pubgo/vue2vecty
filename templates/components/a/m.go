package a

import (
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"github.com/pubgo/g/xerror"
)

func M(data map[string]interface{}, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	t := &_M{Slot: slots}
	decoder := xerror.PanicErr(mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "prop",
		Metadata:         nil,
		Result:           t,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	})).(*mapstructure.Decoder)
	xerror.Panic(decoder.Decode(data))
	return t
}

type _M struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Slot vecty.List
}
