package a

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"github.com/pubgo/g/xerror"
)

func M(data js.M, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
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

	//"github.com/vincent-petithory/dataurl"
	//vecty.AddStylesheet(dataurl.New([]byte(styles), "text/css").String())
	return t
}

type _M struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Slot vecty.List
}
