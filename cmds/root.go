package cmds

import (
	"github.com/pubgo/g/xcmds"
	"github.com/pubgo/g/xerror"
	"github.com/pubgo/vue2vecty/cmds/v2v"
	"github.com/pubgo/vue2vecty/version"
)

const Service = "vue2vecty"
const EnvPrefix = "V2V"

// Execute exec
var Execute = xcmds.Init(EnvPrefix, func(cmd *xcmds.Command) {
	defer xerror.Assert()

	cmd.Use = Service
	cmd.Version = version.Version

	v2v.Version()
	v2v.Build()
	v2v.Component()
})
