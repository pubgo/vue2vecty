// +build !js !wasm

package main

import (
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/cmds"
)

func main() {
	xerror.Exit(cmds.Execute(), "cmd error")
}
