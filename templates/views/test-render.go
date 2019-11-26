// +build js

package views

import (
	"fmt"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"honnef.co/go/js/dom/v2"
)

func init() {
	dom.GetWindow().AddEventListener("test-ssss", false, func(e dom.Event) {
		fmt.Println(e.Target().TextContent())
	})
}

func (t *Test) Render() vecty.ComponentOrHTML {
	return elem.Body(
		elem.Heading1(
			vecty.Text("html2vecty"),
		),
		elem.Paragraph(
			vecty.Text("Enter HTML here and the vecty syntax will appear opposite."),
		),
		func() (e vecty.List) {
			for k := range t.Data {
				e = append(e, elem.Heading2(vecty.Text("Class attributes"), vecty.Markup(
					vecty.Style("border", t.Data[k]),
					vecty.Style("color", "red!important"),
				)))
			}
			return
		}(),
		vecty.If(len("") > 0, elem.Heading2(vecty.Text("Class attributes"))),
		&Test{},
		elem.Paragraph(
			vecty.Markup(
				vecty.Class("foo", "bar", "baz"),
				vecty.MarkupIf(len("") > 0, vecty.Class("foo", "bar", "baz")),
				event.Click(func(i *vecty.Event) {
					dom.GetWindow().DispatchEvent(dom.WrapEvent(i.Target))
					dom.WrapEvent(i.Target).PreventDefault()

				}),
			),
		),
		elem.Heading2(
			vecty.Text("Style attributes"),
		),
		elem.Paragraph(
			vecty.Markup(
				vecty.Style("border", "2px"),
				vecty.Style("color", "red!important"),
			),
		),
		elem.Heading2(
			vecty.Text("Special properties"),
		),
		elem.Input(
			vecty.Markup(
				prop.Type(prop.TypeCheckbox),
				prop.Checked(true),
				prop.Autofocus(true),
			),
		),
		elem.Anchor(
			vecty.Markup(
				prop.Href("href"),
				prop.ID("id"),
				vecty.Data("foo", "bar"),
			),
			vecty.Text("Props"),
		),
		elem.Heading2(
			vecty.Text("An example"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("modal"),
				vecty.Attribute("tabindex", "-1"),
				vecty.Attribute("role", "dialog"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("modal-dialog"),
					vecty.Attribute("role", "document"),
				),
				elem.Div(
					vecty.Markup(
						vecty.Class("modal-content"),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("modal-header"),
						),
						elem.Heading5(
							vecty.Markup(
								vecty.Class("modal-title"),
							),
							vecty.Text("Modal title"),
						),
						elem.Button(
							vecty.Markup(
								prop.Type(prop.TypeButton),
								vecty.Class("close"),
								vecty.Data("dismiss", "modal"),
								vecty.Attribute("aria-label", "Close"),
							),
							elem.Span(
								vecty.Markup(
									vecty.Attribute("aria-hidden", "true"),
								),
								vecty.Text("x"),
							),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("modal-body"),
						),
						elem.Paragraph(
							vecty.Text("Modal body text goes here."),
						),
					),
					elem.Div(
						vecty.Markup(
							vecty.Class("modal-footer"),
						),
						elem.Button(
							vecty.Markup(
								prop.Type(prop.TypeButton),
								vecty.Class("btn", "btn-primary"),
							),
							vecty.Text("Save changes"),
						),
						elem.Button(
							vecty.Markup(
								prop.Type(prop.TypeButton),
								vecty.Class("btn", "btn-secondary"),
								vecty.Data("dismiss", "modal"),
							),
							vecty.Text("Close"),
						),
					),
				),
			),
		),
	)
}
