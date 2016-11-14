package cloudconfig

import (
	"net/rpc"
	"os"

	plugin "github.com/hashicorp/go-plugin"
)

func NewCloudConfigPlugin(plg Deployer) Plugin {
	return Plugin{
		Plugin: plg,
	}
}

type Plugin struct {
	Plugin Deployer
}

func (s Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: s.Plugin}, nil
}

func (s Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPC{client: c}, nil
}

// PluginsMapHash is an identifier for plugins registered with the go-plugin library.
const PluginsMapHash = "cloudconfig"

// HandshakeConfig is the configuration for establishing communication between the CLI plugins.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// Run runs a CloudConfigDeployer as an RPC server.
// It should be called from a plugin's func main.
func Run(cc Deployer) {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				PluginsMapHash: NewCloudConfigPlugin(cc),
			},
		})
		return
	}
}
