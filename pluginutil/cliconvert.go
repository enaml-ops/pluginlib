package pluginutil

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/pluginlib/pcli"
)

//ToCliFlagArray - converts a plugin flag array into a
//codegangsta cli.Flag array
func ToCliFlagArray(fs []pcli.Flag) (cliFlags []cli.Flag) {
	for _, f := range fs {
		cliFlags = append(cliFlags, f.ToCli().(cli.Flag))
	}
	return
}
