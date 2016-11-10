package product

import (
	"errors"
	"net/rpc"

	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
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

// RPC is an implementation of Deployer that talks over RPC.
type RPC struct {
	client *rpc.Client
}

// GetProduct calls a plugin's GetProduct method over RPC.
func (p *RPC) GetProduct(args []string, cloudConfig []byte, cs cred.Store) ([]byte, error) {
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

// GetMeta calls a plugin's GetMeta method over RPC.
func (p *RPC) GetMeta() Meta {
	var resp Meta
	if err := p.client.Call("Plugin.GetMeta", new(interface{}), &resp); err != nil {
		panic(err)
	}
	return resp
}

// GetFlags calls a plugin's GetFlags method over RPC.
func (p *RPC) GetFlags() []pcli.Flag {
	var resp []pcli.Flag
	if err := p.client.Call("Plugin.GetFlags", new(interface{}), &resp); err != nil {
		panic(err)
	}
	return resp
}

// RPCServer is the RPC server that ProductRPC connects to.
// It conforms to the requirements of net/rpc.
type RPCServer struct {
	Impl Deployer
}

// GetProduct forwards the RPC request to the plugin's GetProduct method
// and sends back the results.
func (p *RPCServer) GetProduct(args Args, resp *Response) error {
	var err error
	resp.Bytes, err = p.Impl.GetProduct(args.Args, args.CloudConfig, args.CredStore)

	if err != nil {
		resp.ErrRes = err.Error()
		return err
	}

	resp.ErrRes = ""
	return nil
}

// GetMeta forwards the RPC request to the plugin's GetMeta method
// and sends back the results.
func (s *RPCServer) GetMeta(args interface{}, resp *Meta) error {
	*resp = s.Impl.GetMeta()
	return nil
}

// GetFlags forwards the RPC request to the plugin's GetFlags method
// and sends back the results.
func (s *RPCServer) GetFlags(args interface{}, resp *[]pcli.Flag) error {
	*resp = s.Impl.GetFlags()
	return nil
}
