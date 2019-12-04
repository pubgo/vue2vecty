package views

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"log"
)

func Test(data js.M, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	t := &_Test{Slot: slots}

	if data != nil {
		if err := mapstructure.Decode(data, t); err != nil {
			log.Fatalf("%#v", err)
		}
	}

	return t
}

type _Test struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Slot vecty.List
}
