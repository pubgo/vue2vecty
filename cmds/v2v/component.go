package v2v

import (
	"github.com/pubgo/g/pkg/fileutil"
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func Component() {
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
			_file := filepath.Join(templateHome, "components", args[0])

			if len(names) == 1 {
				xerror.Panic(vue2vecty.CreateComponent("components", names[0]).Save(_file + ".go"))
			} else {
				xerror.Panic(fileutil.IsNotExistMkDir(filepath.Dir(_file + ".go")))
				xerror.Panic(vue2vecty.CreateComponent(names[len(names)-2], names[len(names)-1]).Save(_file + ".go"))
			}

			xerror.Panic(ioutil.WriteFile(_file+".html", []byte(`<p>{{0==0?"world":"hello"}}</p>`), 0644))
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
				xerror.Panic(vue2vecty.CreateComponent("views", names[0]).Save(_file))
			} else {
				xerror.Panic(fileutil.IsNotExistMkDir(filepath.Dir(_file)))
				xerror.Panic(vue2vecty.CreateComponent(names[len(names)-2], names[len(names)-1]).Save(_file))
			}
			xerror.Panic(ioutil.WriteFile(_file+".html", []byte(`<p>{{0==0?"world":"hello"}}</p>`), 0644))
			return
		}},
	))
}
