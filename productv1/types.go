// Package productv1 is the API for the V1 product interface.
package productv1

import (
	"github.com/enaml-ops/pluginlib/cred"
)

// ProductDeployer is the interface implemented by V1 product plugins.
type ProductDeployer interface {
	GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error)
}
