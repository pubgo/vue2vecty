// This file was created with https://github.com/pubgo/vue2vecty
package button

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/mitchellh/mapstructure"
	"log"
)

func ButtonGroup(data js.M, slots ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	t := &_ButtonGroup{Slot: slots}
	if data != nil {
		if err := mapstructure.Decode(data, t); err != nil {
			log.Fatalf("%#v", err)
		}
	}
	return t
}

type _ButtonGroup struct {
	vecty.Core
	Slot vecty.List
}

func (t *_ButtonGroup) Render() vecty.ComponentOrHTML {
	vecty.SetTitle(t.GetTitle())
	return t._Render()
}
