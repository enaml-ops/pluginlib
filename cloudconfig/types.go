package cloudconfig

import (
	"log"
	"net/rpc"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/go-plugin"
)

type Meta struct {
	Name       string
	Properties map[string]interface{}
}

// CloudConfigDeployer is the interface that we will expose for cloud config
// plugins
type CloudConfigDeployer interface {
	GetMeta() Meta
	GetFlags() []cli.Flag
	GetCloudConfig(args []string) []byte
}

// CloudConfigRPC - Here is an implementation that talks over RPC
type CloudConfigRPC struct{ client *rpc.Client }

func (s *CloudConfigRPC) GetMeta() Meta {
	var resp Meta
	err := s.client.Call("Plugin.GetMeta", new(interface{}), &resp)

	if err != nil {
		panic(err)
	}
	return resp
}

func (s *CloudConfigRPC) GetCloudConfig(args []string) []byte {
	var resp []byte
	log.Println("calling rpc client getcloudconfig")
	err := s.client.Call("Plugin.GetCloudConfig", args, &resp)
	log.Println("call:", err)
	if err != nil {
		panic(err)
	}
	return resp
}

func (s *CloudConfigRPC) GetFlags() []cli.Flag {
	var resp []cli.Flag
	err := s.client.Call("Plugin.GetFlags", new(interface{}), &resp)
	log.Println("call: ", err)

	if err != nil {
		panic(err)
	}
	return resp
}

//CloudConfigRPCServer - Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type CloudConfigRPCServer struct {
	Impl CloudConfigDeployer
}

func (s *CloudConfigRPCServer) GetFlags(args interface{}, resp *[]cli.Flag) error {
	*resp = s.Impl.GetFlags()
	return nil
}

func (s *CloudConfigRPCServer) GetMeta(args interface{}, resp *Meta) error {
	*resp = s.Impl.GetMeta()
	return nil
}

func (s *CloudConfigRPCServer) GetCloudConfig(args []string, resp *[]byte) error {
	*resp = s.Impl.GetCloudConfig(args)
	return nil
}

func NewCloudConfigPlugin(plg CloudConfigDeployer) CloudConfigPlugin {
	return CloudConfigPlugin{
		Plugin: plg,
	}
}

type CloudConfigPlugin struct {
	Plugin CloudConfigDeployer
}

func (s CloudConfigPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &CloudConfigRPCServer{Impl: s.Plugin}, nil
}

func (s CloudConfigPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &CloudConfigRPC{client: c}, nil
}
