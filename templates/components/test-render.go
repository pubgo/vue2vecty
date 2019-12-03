// This file was created with https://github.com/pubgo/vue2vecty
// using https://jsgo.io/pubgo/vue2vecty
package components

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	dom "honnef.co/go/js/dom/v2"
)

func (t *_Test) Render() vecty.ComponentOrHTML {
	return elem.Div(elem.Div(elem.Navigation(vecty.Markup(vecty.Class("navbar", "navbar-expand-md", "navbar-dark", "bg-dark", "fixed-top")), elem.Anchor(vecty.Markup(vecty.Class("navbar-brand")), vecty.Text("Navbar")), elem.Button(vecty.Markup(vecty.Class("navbar-toggler"), vecty.Data("toggle", "toggle"), vecty.Data("target", "target")), elem.Span(vecty.Markup(vecty.Class("navbar-toggler-icon")))), elem.Div(vecty.Markup(vecty.Class("collapse", "navbar-collapse")), elem.UnorderedList(vecty.Markup(vecty.Class("navbar-nav", "mr-auto")), elem.ListItem(vecty.Markup(vecty.Class("nav-item", "active")), elem.Anchor(vecty.Markup(vecty.Class("nav-link")), vecty.Text("Home"), elem.Span(vecty.Markup(vecty.Class("sr-only")), vecty.Text("(current)")))), elem.ListItem(vecty.Markup(vecty.Class("nav-item")), elem.Anchor(vecty.Markup(vecty.Class("nav-link")), vecty.Text("Link"))), elem.ListItem(vecty.Markup(vecty.Class("nav-item")), elem.Anchor(vecty.Markup(vecty.Class("nav-link", "disabled")), vecty.Text("Disabled"))), elem.ListItem(vecty.Markup(vecty.Class("nav-item", "dropdown")), elem.Anchor(vecty.Markup(vecty.Class("nav-link", "dropdown-toggle"), vecty.Data("toggle", "toggle")), vecty.Text("Dropdown")), elem.Div(vecty.Markup(vecty.Class("dropdown-menu")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item")), vecty.Text("Action")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item")), vecty.Text("Another action")), elem.Anchor(vecty.Markup(vecty.Class("dropdown-item")), vecty.Text("Something else here"))))), elem.Form(vecty.Markup(vecty.Class("form-inline", "my-2", "my-lg-0")), elem.Input(vecty.Markup(vecty.Class("form-control", "mr-sm-2"))), elem.Button(vecty.Markup(vecty.Class("btn", "btn-outline-success", "my-2", "my-sm-0")), vecty.Text("Search"))))), elem.Div(vecty.Markup(), elem.OrderedList(func() (e vecty.List) {
		for todo := range t.todos {
			e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(todo.text)))
		}
		return
	}())), elem.Div(vecty.Markup(vecty.Style("float", right)), elem.Label(elem.TextArea(vecty.Markup(vecty.Style("font-family", monospace), event.Input(t.texthandler)), vecty.Text("{vecty-field:Input}")))), elem.Div(vecty.Markup(), vecty.Text(message)), elem.Div(vecty.Markup(), elem.Span(vecty.Markup(vecty.Data("title", message)), vecty.Text("鼠标悬停几秒钟查看此处动态绑定的提示信息！"))), elem.Div(vecty.Markup(), vecty.If(seen, elem.Paragraph(vecty.Markup(), vecty.Text("现在你看到我了")))), elem.Div(vecty.Markup(), elem.Paragraph(vecty.Text(message)), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息")), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息")), elem.Div(vecty.Markup(), elem.Paragraph(vecty.Text(message)), elem.Button(vecty.Markup(event.Click(t.reverseMessage)), vecty.Text("反转消息"))), elem.Div(vecty.Markup(), elem.Paragraph(vecty.Text(message)), elem.Input(vecty.Markup(prop.Value(t.message), event.Input(func(e *vecty.Event) {
		t.message = dom.WrapEvent(e.Target).Target().TextContent()
		dom.WrapEvent(e.Target).PreventDefault()
	})))), elem.Div(vecty.Markup(), elem.OrderedList(func() (e vecty.List) {
		for item := range t.groceryList {
			e = append(e, TodoItem(map[string]interface{}{
				"Key":  item.id,
				"Todo": item,
			}, elem.Div(vecty.Markup(), AppNav(), AppView(AppSidebar(), AppContent()))))
		}
		return
	}()))), elem.Div(vecty.Markup(), elem.Paragraph(vecty.Text(message)), elem.Input(vecty.Markup(prop.Value(t.message), event.Input(func(e *vecty.Event) {
		t.message = dom.WrapEvent(e.Target).Target().TextContent()
		dom.WrapEvent(e.Target).PreventDefault()
	})))), elem.OrderedList(TodoItem()), elem.Div(vecty.Markup(), elem.OrderedList(func() (e vecty.List) {
		for item := range t.groceryList {
			e = append(e, TodoItem(map[string]interface{}{
				"Key":  item.id,
				"Todo": item,
			}))
		}
		return
	}())), elem.Div(vecty.Markup(), AppNav(), AppView(AppSidebar(), AppContent()))))
}
