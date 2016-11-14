package product

import (
	"net/rpc"
	"os"

	plugin "github.com/hashicorp/go-plugin"
)

// Plugin wraps up the RPC server and client into a single type.
type Plugin struct {
	Plugin Deployer
}

// Server returns an RPC server that implements the ProductDeployer interface.
func (p Plugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Plugin}, nil
}

// Client returns an RPC client that implements the ProductDeployer interface.
func (p Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPC{client: c}, nil
}

// NewProductPlugin decorates a ProductDeployer with the RPC functionality
// requried to operate as a product plugin.
func NewProductPlugin(pd Deployer) Plugin {
	return Plugin{Plugin: pd}
}

// PluginsMapHash is an identifier for plugins registered with the go-plugin library.
const PluginsMapHash = "product"

// HandshakeConfig is the configuration for establishing communication between the CLI plugins.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// Run runs a ProductDeployer as an RPC server.
// It should be called from a plugin's func main.
func Run(p Deployer) {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				PluginsMapHash: NewProductPlugin(p),
			},
		})
		return
	}
}
