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
			),
			vecty.Text(""),
		),
		vecty.Text(""),
	)
}

func TestC(data map[string]interface{}, slot ...vecty.ComponentOrHTML) vecty.ComponentOrHTML {
	return &Test{}
}

type Test struct {
	vecty.Core

	OnInput func(string)

	hello string
	Value string

	Slot   vecty.List
	Markup []vecty.Applyer
	data   map[string]interface{}
}
