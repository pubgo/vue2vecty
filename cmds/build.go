package cmds

import (
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"github.com/spf13/cobra"
)

func init() {
	xcmds.AddCommand(&cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build template",
		Run: func(cmd *cobra.Command, args []string) {
			vue2vecty.NewTranspiler()
		},
	})
}
