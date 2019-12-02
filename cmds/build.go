package cmds

import (
	"fmt"
	"github.com/gobuffalo/envy"
	"github.com/pubgo/g/errors"
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/vue2vecty"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	var templateHome = "templates"
	xcmds.AddCommand(&xcmds.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build template",
		Run: func(cmd *xcmds.Command, args []string) {
			defer xerror.Assert()

			xerror.Panic(filepath.Walk(templateHome, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || !strings.Contains(info.Name(), ".") {
					return xerror.Wrap(err, "file walk failed")
				}

				names := strings.Split(info.Name(), ".")
				if !strings.Contains(names[len(names)-1], "html") {
					return nil
				}

				_dir := filepath.Dir(path)
				_compo := filepath.Base(_dir)
				name := names[0]
				fmt.Println(_dir, _compo, name)

				f, err := os.Open(path)
				xerror.PanicM(err, "file open error")

				_c := vue2vecty.NewTranspiler(f, envy.CurrentPackage()+"/templates", strings.ReplaceAll(strings.Title(name), "-", ""), _compo)
				errors.Panic(ioutil.WriteFile(filepath.Join(_dir, name+"-render.go"), []byte(_c.Code()), 0755))
				return nil
			}))

		},
	})
}
