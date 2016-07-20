package pcli

import "encoding/gob"

func init() {
	gob.Register(Flag{})
}
