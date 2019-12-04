// This file was created with https://github.com/pubgo/vue2vecty. DO NOT EDIT.
// using https://jsgo.io/pubgo/vue2vecty
package views

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"github.com/pubgo/vue2vecty/templates/components"
	dom "honnef.co/go/js/dom/v2"
)

func (t *_Test) _Render() vecty.ComponentOrHTML {
	vecty.SetTitle(t.GetTitle())
	return elem.Body(elem.Div(elem.Div(elem.Navigation(vecty.Markup(vecty.Class("navbar", "navbar-expand-md", "navbar-dark", "bg-dark", "fixed-top")), elem.Anchor(vecty.Markup(vecty.Class("navbar-brand"), vecty.Property("href", "#")), vecty.Text("Navbar")), elem.Button(vecty.Markup(vecty.Class("navbar-toggler"), vecty.Property("type", "button"), vecty.Data("toggle", "toggle"), vecty.Data("target", "target"), vecty.Property("aria-controls", "navbarsExampleDefault"), vecty.Property("aria-expanded", "false"), vecty.Property("aria-label", "Toggle navigation")), elem.Span(vecty.Markup(vecty.Class("navbar-toggler-icon")))), elem.Div(vecty.Markup(vecty.Class("collapse", "navbar-collapse"), vecty.Property("id", "navbarsExampleDefault")), elem.UnorderedList(vecty.Markup(vecty.Class("navbar-nav", "mr-auto")), elem.ListItem(vecty.Markup(vecty.Class("nav-item", "active")), elem.Anchor(vecty.Markup(vecty.Class("nav-link"), vecty.Property("href", "#")), vecty.Text("Home"), elem.Span(vecty.Markup(vecty.Class("sr-only")), vecty.Text("(current)")))), elem.ListItem(vecty.Markup(vecty.Class("nav-item")), elem.Anchor(vecty.Markup(vecty.Class("nav-link"), vecty.Property("href", "#")), vecty.Text("Link"))), elem.ListItem(vecty.Markup(vecty.Class("nav-item")), elem.Anchor(vecty.Markup(vecty.Class("nav-link", "disabled"), vecty.Property("href", "#")), vecty.Text("Disabled"))), elem.ListItem(vecty.Markup(vecty.Class("nav-item", "dropdown")), elem.Anchor(vecty.Markup(vecty.Class("nav-link", "dropdown-toggle"), vecty.Property("href", "https://example.com"), vecty.Property("id", "dropdown01"), vecty.Data("toggle", "toggle"), vecty.Property("aria-haspopup", "true"), vecty.Property("aria-expanded", "false")), vecty.Text("Dropdown")), elem.Div(vecty.Markup(vecty.Class("dropdown-menu"), vecty.Property("aria-labelledby", "dropdown01")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item"), vecty.Property("href", "#")), vecty.Text("Action")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item"), vecty.Property("href", "#")), vecty.Text("Another action")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item"), vecty.Property("href", "#")), vecty.Text("Something else here"))))), elem.Form(vecty.Markup(vecty.Class("form-inline", "my-2", "my-lg-0")), elem.Input(vecty.Markup(vecty.Class("form-control", "mr-sm-2"), vecty.Property("type", "text"), vecty.Property("placeholder", "Search"), vecty.Property("aria-label", "Search"))), elem.Button(vecty.Markup(vecty.Class("btn", "btn-outline-success", "my-2", "my-sm-0"), vecty.Property("type", "submit")), vecty.Text("Search"))))), elem.Div(vecty.Markup(vecty.Property("id", "app-4")), elem.OrderedList(func() (e vecty.List) {
		for __todo := range t.todos {
			todo := __todo
			e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(todo.text)))
		}
		return
	}())), elem.Div(vecty.Markup(vecty.Style("float", right)), elem.Label(elem.TextArea(vecty.Markup(vecty.Style("font-family", monospace), vecty.Property("cols", "70"), vecty.Property("rows", "14"), event.Input(t.texthandler)), vecty.Text("{vecty-field:Input}")))), elem.Div(vecty.Markup(vecty.Property("id", "app")), vecty.Text(message)), elem.Div(vecty.Markup(vecty.Property("id", "app-2")), elem.Span(vecty.Markup(vecty.Property("title", message)), vecty.Text("鼠标悬停几秒钟查看此处动态绑定的提示信息！"))), elem.Div(vecty.Markup(vecty.Property("id", "app-3")), vecty.If(seen, elem.Paragraph(vecty.Markup(), vecty.Text("现在你看到我了")))), elem.Div(vecty.Markup(vecty.Property("id", "app-5")), elem.Paragraph(vecty.Text(message)), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息")), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息")), elem.Div(vecty.Markup(vecty.Property("id", "app-5")), elem.Paragraph(vecty.Text(message)), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息"))), elem.Div(vecty.Markup(vecty.Property("id", "app-6")), elem.Paragraph(vecty.Text(message)), elem.Input(vecty.Markup(prop.Value(t.message), event.Input(func(e *vecty.Event) {
		t.message = dom.WrapEvent(e.Target).Target().TextContent()
		dom.WrapEvent(e.Target).PreventDefault()
	})))), elem.Div(vecty.Markup(vecty.Property("id", "app-7")), elem.OrderedList(func() (e vecty.List) {
		for __item := range t.groceryList {
			item := __item
			e = append(e, components.TodoItem(map[string]interface{}{
				"Key":  item.id,
				"Todo": item,
			}, elem.Div(vecty.Markup(vecty.Property("id", "app")), components.AppNav(), components.AppView(components.AppSidebar(), components.AppContent()))))
		}
		return
	}()))), elem.Div(vecty.Markup(vecty.Property("id", "app-6")), elem.Paragraph(vecty.Text(message)), elem.Input(vecty.Markup(prop.Value(t.message), event.Input(func(e *vecty.Event) {
		t.message = dom.WrapEvent(e.Target).Target().TextContent()
		dom.WrapEvent(e.Target).PreventDefault()
	})))), elem.OrderedList(components.TodoItem()), elem.Div(vecty.Markup(vecty.Property("id", "app-7")), elem.OrderedList(func() (e vecty.List) {
		for __item := range t.groceryList {
			item := __item
			e = append(e, components.TodoItem(map[string]interface{}{
				"Key":  item.id,
				"Todo": item,
			}))
		}
		return
	}())), elem.Div(vecty.Markup(vecty.Property("id", "app")), components.AppNav(), components.AppView(components.AppSidebar(), components.AppContent())))))
}
