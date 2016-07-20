package pluginutil

import (
	"github.com/codegangsta/cli"
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
