package cloudconfig

import "github.com/enaml-ops/pluginlib/pcli"

type Meta struct {
	Name       string
	Properties map[string]interface{}
}

// CloudConfigDeployer is the interface that we will expose for cloud config
// plugins
type CloudConfigDeployer interface {
	GetMeta() Meta
	GetFlags() []pcli.Flag
	GetCloudConfig(args []string) ([]byte, error)
}
