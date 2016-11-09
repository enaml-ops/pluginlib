package cloudconfig

import "github.com/enaml-ops/pluginlib/pcli"

type Meta struct {
	Name       string
	Properties map[string]interface{}
}

// Deployer is the interface for cloud config plugins
type Deployer interface {
	GetMeta() Meta
	GetFlags() []pcli.Flag
	GetCloudConfig(args []string) ([]byte, error)
}
