// This file was created with https://github.com/pubgo/vue2vecty. DO NOT EDIT.
// using https://jsgo.io/pubgo/vue2vecty
// +build js wasm

package button

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/pubgo/vue2vecty/templates/components"
)

func (t *_ButtonGroup) _Render() vecty.ComponentOrHTML {
	return elem.Paragraph(vecty.Text(func() string {
			if 0 == 0 {
				return "world"
			} else {
				return "hello"
			}
		}())), elem.UnorderedList(vecty.Markup(vecty.Property("id", "example-1")), func() (e vecty.List) {
			for __item := range t.items {
				item := __item
				e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(item.message)))
			}
			return
		}(), func() (e vecty.List) {
			for __item, __index := range t.items {
				item, index := __item, __index
				e = append(e, elem.ListItem(vecty.Markup(), vecty.Text("{{ parentMessage }} - {{ index }} - "+item.message)))
			}
			return
		}(), func() (e vecty.List) {
			for __value, __name := range t.object {
				value, name := __value, __name
				e = append(e, elem.Div(vecty.Markup(), vecty.Text("{{ index }}. {{ name }}: "+value)))
			}
			return
		}(), func() (e vecty.List) {
			for __item := range t.items {
				item := __item
				e = append(e, elem.Div(vecty.Markup(vecty.Property("key", item.id))))
			}
			return
		}(), func() (e vecty.List) {
			for __n := range t.evenNumbers {
				n := __n
				e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(n)))
			}
			return
		}(), vecty.If(!todo.isComplete, func() (e vecty.List) {
			for __todo := range t.todos {
				todo := __todo
				e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(todo)))
			}
			return
		}()), vecty.If(todos.length, elem.UnorderedList(vecty.Markup(), func() (e vecty.List) {
			for __todo := range t.todos {
				todo := __todo
				e = append(e, elem.ListItem(vecty.Markup(), vecty.Text(todo)))
			}
			return
		}())), func() (e vecty.List) {
			for __item := range t.items {
				item := __item
				e = append(e, components.MyComponent(js.M{"Key": item.id}))
			}
			return
		}(), func() (e vecty.List) {
			for __item, __index := range t.items {
				item, index := __item, __index
				e = append(e, components.MyComponent(js.M{
					"Index": index,
					"Item":  item,
					"Key":   item.id,
				}))
			}
			return
		}())
}
