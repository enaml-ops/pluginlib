package pcli

import "encoding/gob"

func init() {
	gob.Register(StringSliceFlag{})
	gob.Register(StringFlag{})
	gob.Register(BoolFlag{})
	gob.Register(BoolTFlag{})
	gob.Register(IntFlag{})
}
