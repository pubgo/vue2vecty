// This file was created with https://github.com/pubgo/vue2vecty
package components

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"log"
)

type _HeeKkk struct {
	vecty.Core
	Slot vecty.List
}

func HeeKkk(data js.M, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	t := &_HeeKkk{Slot: slots}
	if data != nil {
		if err := mapstructure.Decode(data, t); err != nil {
			log.Fatalf("%#v", err)
		}
	}
	return t
}
