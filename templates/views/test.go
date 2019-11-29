package views

import (
	"github.com/gopherjs/vecty"
)

func init() {
	TestC(map[string]interface{}{
		"OnInput": func(string) {},
		"Value":   "",
	},
		TestC(map[string]interface{}{
			"OnInput": func(string) {},
			"Value":   "",
		},
			TestC(map[string]interface{}{
				"OnInput": func(string) {},
				"Value":   "",
			},
				vecty.Text(""),
				vecty.Markup(),
			),
			vecty.Text(""),
			vecty.Markup(),
		),
		vecty.Text(""),
		vecty.Markup(),
	)
}

func TestC(data map[string]interface{}, child ...vecty.MarkupOrChild) vecty.ComponentOrHTML {
	return &Test{}
}

type Test struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Markup []vecty.Applyer
	Slot   []vecty.ComponentOrHTML
	data   map[string]interface{}
}
