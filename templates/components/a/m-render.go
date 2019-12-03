// This file was created with https://github.com/pubgo/vue2vecty
// using https://jsgo.io/pubgo/vue2vecty
package a

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/pubgo/vue2vecty/templates/components/b/a"
	"github.com/pubgo/vue2vecty/templates/components/c/a/b"
)

func (t *_M) Render() vecty.ComponentOrHTML {
	return func() (e vecty.List) {
		for key, value := range t.groceryList {
			e = append(e, a.TodoItem(map[string]interface{}{
				"Key":    item.id,
				"Key3":   item.id,
				"OnKey4": item.id,
				"OnKey5": item.id,
				"Todo":   item,
				t.key1:   item.id,
			}, vecty.Text("sss"), func() (e vecty.List) {
				for key, value := range t.groceryList {
					e = append(e, b.D(map[string]interface{}{
						"Hello":   "hello",
						"Key":     key.id,
						"OnClick": t.onClick,
						"OnInput": func(v string) {
							t.todo = v
						},
						"OnKey": t.click,
						"Value": t.todo,
					}))
				}
				return
			}(), elem.ListItem(), elem.Input(), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text("0?world:\"hello\"")), elem.Paragraph(vecty.Markup(vecty.Data("click-sss", "click-sss")), vecty.Text(func() string {
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
