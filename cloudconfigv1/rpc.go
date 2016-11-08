package cloudconfig

import (
	"errors"
	"log"
	"net/rpc"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
)

// Response contains the results of a GetCloudConfig RPC call.
type Response struct {
	Bytes  []byte
	ErrRes string
}

// RPC - Here is an implementation that talks over RPC
type RPC struct{ client *rpc.Client }

func (s *RPC) GetMeta() Meta {
	var resp Meta
	err := s.client.Call("Plugin.GetMeta", new(interface{}), &resp)

	if err != nil {
		log.Println("[ERROR] GetFlags: ", err)
	}
	return resp
}

func (s *RPC) GetCloudConfig(args []string) ([]byte, error) {
	var resp Response
	lo.G.Debug("calling rpc client getcloudconfig")
	err := s.client.Call("Plugin.GetCloudConfig", args, &resp)
	if err != nil {
		lo.G.Debug("[ERROR] GetCloudConfig:", err)
		return nil, err
	}
	if resp.ErrRes != "" {
		lo.G.Debug("error:", resp.ErrRes)
		return nil, errors.New(resp.ErrRes)
	}
	return resp.Bytes, nil
}

func (s *RPC) GetFlags() []pcli.Flag {
	var resp []pcli.Flag
	err := s.client.Call("Plugin.GetFlags", new(interface{}), &resp)

	if err != nil {
		log.Println("[ERROR] GetFlags: ", err)
		return nil
	}
	return resp
}

// RPCServer - Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	Impl CloudConfigDeployer
}

func (s *RPCServer) GetFlags(args interface{}, resp *[]pcli.Flag) error {
	*resp = s.Impl.GetFlags()
	return nil
}

func (s *RPCServer) GetMeta(args interface{}, resp *Meta) error {
	*resp = s.Impl.GetMeta()
	return nil
}

func (s *RPCServer) GetCloudConfig(args []string, resp *Response) error {
	var err error
	resp.Bytes, err = s.Impl.GetCloudConfig(args)

	if err != nil {
		resp.ErrRes = err.Error()
		return err
	}

	resp.ErrRes = ""
	return nil
}
