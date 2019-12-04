// This file was created with https://github.com/pubgo/vue2vecty. DO NOT EDIT.
// using https://jsgo.io/pubgo/vue2vecty
package a

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/pubgo/vue2vecty/templates/components"
)

func (t *_M) _Render() vecty.ComponentOrHTML {
	return func() (e vecty.List) {
		for __key, __value := range t.groceryList {
			key, value := __key, __value
			e = append(e, components.TodoItem(map[string]interface{}{
				".Key2":  "item.id",
				"Key":    item.id,
				"Key3":   item.id,
				"OnKey4": t.item.id,
				"OnKey5": t.Panic,
				"Todo":   item,
				t.key1:   item.id,
			}, vecty.Text("sss"), vecty.If(a > 1 && b == 2 && key["sss"]+1 > 0, func() (e vecty.List) {
				for __k := range t.groceryList() {
					k := __k
					e = append(e, components.Cd(map[string]interface{}{
						"Hello":   "hello",
						"Key":     key.id,
						"OnClick": t.onClick,
						"OnInput": func(_value string) {
							t.todo = _value
							vecty.Rerender(t)
						},
						"OnKey": t.click(k),
						"OnKey2": func(value string) {
							t.mk = key.id + value
							vecty.Rerender(t)
						},
						"Value": t.todo,
					}))
				}
				return
			}()), elem.ListItem(), elem.Input(), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text("0?world:\"hello\"")), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text("0?world:\"hello\"")), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text(func() string {
				if 0 > 0 {
					return world
				} else {
					return "hello"
				}
			}())), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text(func() string {
				if 0 > 0 {
					return world
				} else {
					return "hello"
				}
			}()))))
		}
		return
	}()
}
