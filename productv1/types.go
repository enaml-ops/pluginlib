// Package product is the API for the V1 product interface.
package product

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
)

// Meta is the metadata for a product plugin.
type Meta struct {
	Name       string
	Properties map[string]interface{}
	Releases   []enaml.Release
	Stemcell   enaml.Stemcell
}

// Deployer is the interface implemented by V1 product plugins.
type Deployer interface {
	GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error)
	GetMeta() Meta
	GetFlags() []pcli.Flag
}
