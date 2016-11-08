package cloudconfig

import (
	"net/rpc"

	plugin "github.com/hashicorp/go-plugin"
)

func NewCloudConfigPlugin(plg CloudConfigDeployer) Plugin {
	return Plugin{
		Plugin: plg,
	}
}

type Plugin struct {
	Plugin CloudConfigDeployer
}

func (s Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: s.Plugin}, nil
}

func (s Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPC{client: c}, nil
}
