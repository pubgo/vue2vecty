// +build js

package views

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"honnef.co/go/js/dom/v2"
)

func init() {
	// <textarea @test-input="hello">
	/*
		<input v-model="searchText">
		等价于
		<input v-bind:value="searchText" v-on:input="searchText = $event.target.value">
		<input v-bind:value="searchText" v-on:input=func(event){
		    searchText = $event.target.value
		}>
		<custom-input v-bind:value="$searchText" v-on:input="searchTextEventFunction" ></custom-input>
		&CustomInput{
			Value=t.searchText,
			OnInput=func(s string) {
				t.searchText = s
			},
		}
	*/
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
		func() vecty.ComponentOrHTML {
			//dom.GetWindow().AddEventListener("test-input-time", false, func(e dom.Event) {
			//	t.hello = e.Target().TextContent()
			//	vecty.Rerender(t)
			//})
			c := &Test{
				OnInput: func(s string) {
					t.Value = s
				},
				Value: t.Value,
			}

			c.Markup = append(c.Markup, vecty.Markup())
			c.Slot = append(c.Slot, vecty.Text(""))

			return c
		}(),
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
				prop.Value(t.Value),
				event.Input(func(event *vecty.Event) {
					t.Value = dom.WrapEvent(event.Target).Target().TextContent()
					dom.WrapEvent(event.Target).PreventDefault()
					// $emit(t.Input, t.Value)
					prop.Value()
				}),
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
