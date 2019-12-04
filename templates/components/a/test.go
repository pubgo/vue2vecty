package a

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"github.com/pubgo/g/xerror"
)

func Test(data js.M, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	t := &_Test{Slot: slots}
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

type _Test struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Slot vecty.List
}
