package product

import (
	"errors"
	"log"
	"net/rpc"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/hashicorp/go-plugin"
)

type Meta struct {
	Name       string
	Properties map[string]interface{}
	Releases   []enaml.Release
	Stemcell   enaml.Stemcell
}

// ProductDeployer is the interface that we will expose for product
// plugins
type ProductDeployer interface {
	GetMeta() Meta
	GetFlags() []pcli.Flag
	GetProduct(args []string, cloudConfig []byte) ([]byte, error)
}

// ProductRPC - Here is an implementation that talks over RPC
type ProductRPC struct{ client *rpc.Client }

func (s *ProductRPC) GetMeta() Meta {
	var resp Meta
	err := s.client.Call("Plugin.GetMeta", new(interface{}), &resp)

	if err != nil {
		panic(err)
	}
	return resp
}

type RPCArgs struct {
	Arg1 []string
	Arg2 []byte
}

type RPCResponse struct {
	ResBytes []byte
	ErrRes   string
}

func (s *ProductRPC) GetProduct(args []string, cloudConfig []byte) (b []byte, err error) {
	var rpcRes = RPCResponse{}
	log.Println("calling rpc client get product")
	if err := s.client.Call("Plugin.GetProduct", RPCArgs{
		Arg1: args,
		Arg2: cloudConfig,
	}, &rpcRes); err != nil {
		log.Println("call:", err)
		return nil, err
	}

	if rpcRes.ErrRes != "" {
		log.Println("error: ", rpcRes.ErrRes)
		err = errors.New(rpcRes.ErrRes)
	}
	return rpcRes.ResBytes, err
}

func (s *ProductRPC) GetFlags() []pcli.Flag {
	var resp []pcli.Flag
	err := s.client.Call("Plugin.GetFlags", new(interface{}), &resp)
	log.Println("call: ", err)

	if err != nil {
		panic(err)
	}
	return resp
}

//ProductRPCServer - Here is the RPC server that ProductRPC talks to, conforming to
// the requirements of net/rpc
type ProductRPCServer struct {
	Impl ProductDeployer
}

func (s *ProductRPCServer) GetFlags(args interface{}, resp *[]pcli.Flag) error {
	*resp = s.Impl.GetFlags()
	return nil
}

func (s *ProductRPCServer) GetMeta(args interface{}, resp *Meta) error {
	*resp = s.Impl.GetMeta()
	return nil
}

func (s *ProductRPCServer) GetProduct(args RPCArgs, resp *RPCResponse) error {
	var err error
	resp.ResBytes, err = s.Impl.GetProduct(args.Arg1, args.Arg2)
	resp.ErrRes = ""

	if err != nil {
		log.Println("we are seeing an error", err)
		resp.ErrRes = err.Error()
	}
	return nil
}

func NewProductPlugin(plg ProductDeployer) ProductPlugin {
	return ProductPlugin{
		Plugin: plg,
	}
}

type ProductPlugin struct {
	Plugin ProductDeployer
}

func (s ProductPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &ProductRPCServer{Impl: s.Plugin}, nil
}

func (s ProductPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ProductRPC{client: c}, nil
}
