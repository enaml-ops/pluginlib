package cloudconfig

import (
	"encoding/gob"

	"gopkg.in/urfave/cli.v2"
)

func init() {
	gob.Register(cli.StringSliceFlag{})
	gob.Register(cli.StringFlag{})
	gob.Register(cli.BoolFlag{})
	gob.Register(cli.BoolFlag{Value: true,})
	gob.Register(cli.DurationFlag{})
	gob.Register(cli.GenericFlag{})
	gob.Register(cli.IntFlag{})
	gob.Register(cli.IntSliceFlag{})
}
