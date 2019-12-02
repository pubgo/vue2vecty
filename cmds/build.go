package cmds

import (
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/vue2vecty/vue2vecty"
)

func init() {
	xcmds.AddCommand(&xcmds.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build template",
		Run: func(cmd *xcmds.Command, args []string) {
			vue2vecty.NewTranspiler()
		},
	})
}
