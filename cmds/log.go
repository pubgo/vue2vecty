package cmds

import (
	"github.com/pubgo/g/xconfig/xconfig_log"
	"github.com/pubgo/g/xinit"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func init() {
	xinit.InitInvoke(func(log *xconfig_log.Log) {
		logger = log.Log("pkg", "vue2vecty")
	})
}
