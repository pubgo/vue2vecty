package cmds

import (
	"github.com/pubgo/g/logs"
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/g/xcmds/xcmd_ss"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/version"
	"github.com/spf13/cobra"
)

const Service = "vue2vecty"
const EnvPrefix = "V2V"

// Execute exec
var Execute = xcmds.Init(EnvPrefix, func(cmd *cobra.Command) {
	defer xerror.Assert()

	cmd.Use = Service
	cmd.Version = version.Version

	// 添加加密命令
	xcmd_ss.Init()
}, func() error {
	_l := logs.Default()
	_l.Version = version.Version
	_l.Service = Service
	_l.Init()

	return nil
})
