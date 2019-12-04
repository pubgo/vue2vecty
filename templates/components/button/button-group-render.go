// This file was created with https://github.com/pubgo/vue2vecty. DO NOT EDIT.
// using https://jsgo.io/pubgo/vue2vecty
package button

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
)

func (t *_ButtonGroup) _Render() vecty.ComponentOrHTML {
	return elem.Paragraph(vecty.Text(func() string {
		if 0 == 0 {
			return "world"
		} else {
			return "hello"
		}
	}()))
}
