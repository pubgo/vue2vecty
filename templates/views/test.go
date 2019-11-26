package views

import (
	"github.com/gopherjs/vecty"
)

type Test struct {
	vecty.Core

	Data []string

	Slot []vecty.MarkupOrChild
}
