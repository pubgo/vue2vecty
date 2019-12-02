package cmds

import (
	"fmt"
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/vue2vecty/version"
	"runtime"
)

func init() {
	xcmds.AddCommand(&xcmds.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "version info",
		Run: func(cmd *xcmds.Command, args []string) {
			fmt.Println("Version:", version.Version)
			fmt.Println("GitHash:", version.CommitV)
			fmt.Println("BuildDate:", version.BuildV)
			fmt.Println("GoVersion:", runtime.Version())
			fmt.Println("GOOS:", runtime.GOOS)
			fmt.Println("GOARCH:", runtime.GOARCH)
		},
	})
}
