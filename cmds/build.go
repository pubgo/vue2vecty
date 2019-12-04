// +build linux windows darwin

package cmds

import (
	"github.com/fsnotify/fsnotify"
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
	var isMonitor = false
	var transpiler = func(path string) (err error) {
		defer xerror.RespErr(&err)
		_dir := filepath.Dir(path)
		_compo := filepath.Base(_dir)
		names := strings.Split(filepath.Base(path), ".")
		name := names[0]

		f, err := os.Open(path)
		xerror.PanicM(err, "file open error")

		_c := vue2vecty.NewTranspiler(f,
			xerror.PanicErr(envy.CurrentModule()).(string)+"/"+templateHome,
			"_"+strings.ReplaceAll(strings.Title(name), "-", ""),
			_compo)
		errors.Panic(ioutil.WriteFile(filepath.Join(_dir, name+"-render.go"), []byte(_c.Code()), 0644))
		return
	}
	xcmds.AddCommand(func(cmd *xcmds.Command) *xcmds.Command {
		cmd.Flags().StringVar(&templateHome, "dir", templateHome, "templates directory")
		cmd.Flags().BoolVarP(&isMonitor, "monitor", "m", isMonitor, "is monitor mode")
		return cmd
	}(&xcmds.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "build template",
		RunE: func(cmd *xcmds.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			watcher := xerror.PanicErr(fsnotify.NewWatcher()).(*fsnotify.Watcher)
			defer watcher.Close()

			xerror.Panic(filepath.Walk(templateHome, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() || !strings.Contains(info.Name(), ".") {
					return xerror.Wrap(err, "file walk failed")
				}

				names := strings.Split(info.Name(), ".")
				// 检查后缀
				if !strings.Contains(names[len(names)-1], "html") {
					return nil
				}

				xerror.Panic(watcher.Add(path))
				xerror.Panic(transpiler(path))
				return nil
			}))

			for isMonitor {
				select {
				case event, ok := <-watcher.Events:
					if ok && (event.Op == fsnotify.Write || event.Op == fsnotify.Rename) {
						_file := event.Name
						logger.Info().Str("file", _file).Str("operation", event.Op.String()).Msg("file change")
						xerror.Panic(transpiler(_file))
						xerror.Panic(watcher.Add(_file))
					}
				case err, ok := <-watcher.Errors:
					if ok {
						xerror.Panic(err)
					}
				}
			}
			return
		}},
	))
}
