package productv1

import (
	"errors"
	"net/rpc"

	"github.com/enaml-ops/pluginlib/cred"
	"github.com/xchapter7x/lo"
)

type (
	// Args contains the args for a GetProduct call.
	Args struct {
		Args        []string
		CloudConfig []byte
		CredStore   cred.Store
	}
	// Response contains the results of a GetProduct call.
	Response struct {
		Bytes  []byte
		ErrRes string
	}
)

// ProductRPC is an implementation of ProductDeployer that talks over RPC.
type ProductRPC struct {
	client *rpc.Client
}

// GetProduct calls a plugin's GetProduct method over RPC.
func (p *ProductRPC) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
	lo.G.Debug("calling RPC client GetProduct")
	var resp Response
	err := p.client.Call("Plugin.GetProduct", Args{
		Args:        args,
		CloudConfig: cloudConfig,
		CredStore:   cs,
	}, &resp)
	if err != nil {
		return nil, err
	}

	if resp.ErrRes != "" {
		lo.G.Debug("error:", resp.ErrRes)
		return nil, errors.New(resp.ErrRes)
	}

	return resp.Bytes, nil
}

// ProductRPCServer is the RPC server that ProductRPC connects to.
// It conforms to the requirements of net/rpc.
type ProductRPCServer struct {
	Impl ProductDeployer
}

// GetProduct forwards the RPC request to the plugin's GetProduct method
// and sends back the results.
func (p *ProductRPCServer) GetProduct(args Args, resp *Response) error {
	var err error
	resp.Bytes, err = p.Impl.GetProduct(args.Args, args.CloudConfig, args.CredStore)

	if err != nil {
		resp.ErrRes = err.Error()
		return err
	}

	resp.ErrRes = ""
	return nil
}
