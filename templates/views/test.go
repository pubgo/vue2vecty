package views

import (
	"github.com/gopherjs/vecty"
	"honnef.co/go/js/dom/v2"
)

type Test struct {
	vecty.Core

	_event dom.Event

	Data []string

	Slot []vecty.MarkupOrChild
}
