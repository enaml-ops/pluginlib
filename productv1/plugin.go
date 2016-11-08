package product

import (
	"net/rpc"
	"os"

	plugin "github.com/hashicorp/go-plugin"
)

// ProductPlugin wraps up the RPC server and client into a single type.
type ProductPlugin struct {
	Plugin ProductDeployer
}

// Server returns an RPC server that implements the ProductDeployer interface.
func (p ProductPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ProductRPCServer{Impl: p.Plugin}, nil
}

// Client returns an RPC client that implements the ProductDeployer interface.
func (p ProductPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProductRPC{client: c}, nil
}

// NewProductPlugin decorates a ProductDeployer with the RPC functionality
// requried to operate as a product plugin.
func NewProductPlugin(pd ProductDeployer) ProductPlugin {
	return ProductPlugin{Plugin: pd}
}

const PluginsMapHash = "product"

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// Run runs a ProductDeployer as an RPC server.
// It should be called from a plugin's func main.
func Run(p ProductDeployer) {
	if len(os.Args) >= 2 && os.Args[1] != "" {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: HandshakeConfig,
			Plugins: map[string]plugin.Plugin{
				"product": NewProductPlugin(p),
			},
		})
	}
}
