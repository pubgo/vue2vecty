package cmds

import (
	"github.com/pubgo/g/errors"
	"github.com/pubgo/g/pkg/fileutil"
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"path/filepath"
	"strings"
)

func init() {
	var templateHome = "templates"

	xcmds.AddCommand(func(cmd *xcmds.Command) *xcmds.Command {
		cmd.Flags().StringVar(&templateHome, "dir", templateHome, "模板目录位置")
		return cmd
	}(&xcmds.Command{
		Use:     "component",
		Aliases: []string{"c", "comp"},
		Short:   "create component",
		RunE: func(cmd *xcmds.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			xerror.PanicT(len(args) == 0, "component name must be specified")
			names := strings.Split(args[0], "/")
			_file := filepath.Join(templateHome, "components", args[0]+".go")

			if len(names) == 1 {
				errors.Panic(vue2vecty.CreateStruct("components", names[0]).Save(_file))
			} else {
				errors.Panic(fileutil.IsNotExistMkDir(filepath.Dir(_file)))
				errors.Panic(vue2vecty.CreateStruct(names[len(names)-2], names[len(names)-1]).Save(_file))
			}
			return
		}},
	))

	xcmds.AddCommand(func(cmd *xcmds.Command) *xcmds.Command {
		cmd.Flags().StringVar(&templateHome, "dir", templateHome, "模板目录位置")
		return cmd
	}(&xcmds.Command{
		Use:     "page",
		Aliases: []string{"p", "view", "route", "r"},
		Short:   "create page",
		RunE: func(cmd *xcmds.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			xerror.PanicT(len(args) == 0, "page name must be specified")
			names := strings.Split(args[0], "/")
			_file := filepath.Join(templateHome, "views", args[0]+".go")
			if len(names) == 1 {
				errors.Panic(vue2vecty.CreateStruct("views", names[0]).Save(_file))
			} else {
				errors.Panic(fileutil.IsNotExistMkDir(filepath.Dir(_file)))
				errors.Panic(vue2vecty.CreateStruct(names[len(names)-2], names[len(names)-1]).Save(_file))
			}
			return
		}},
	))
}
