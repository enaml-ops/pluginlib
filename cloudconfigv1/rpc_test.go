package cloudconfig_test

import (
	"github.com/enaml-ops/pluginlib/cloudconfigv1"
	"github.com/enaml-ops/pluginlib/cloudconfigv1/cloudconfigv1fakes"
	"github.com/enaml-ops/pluginlib/pcli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cloudconfigv1 RPC", func() {
	var d *cloudconfigv1fakes.FakeDeployer

	BeforeEach(func() {
		d = new(cloudconfigv1fakes.FakeDeployer)
	})

	It("Forwards calls to GetMeta", func() {
		controlMeta := cloudconfig.Meta{
			Name: "fakemeta",
		}
		d.GetMetaReturns(controlMeta)
		rpc := cloudconfig.RPCServer{
			Impl: d,
		}

		var resp cloudconfig.Meta
		Ω(rpc.GetMeta(new(interface{}), &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlMeta))
	})

	It("Forwards calls to GetCloudConfig", func() {
		controlResp := cloudconfig.Response{
			Bytes:  []byte{0, 1, 2},
			ErrRes: "",
		}
		d.GetCloudConfigReturns(controlResp.Bytes, nil)
		rpc := cloudconfig.RPCServer{Impl: d}

		var resp cloudconfig.Response
		Ω(rpc.GetCloudConfig([]string{"cc"}, &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlResp))
	})

	It("Forwards calls to GetFlags", func() {
		controlFlags := []pcli.Flag{
			pcli.CreateStringFlag("str", "dummy", ""),
		}
		d.GetFlagsReturns(controlFlags)
		rpc := cloudconfig.RPCServer{
			Impl: d,
		}

		var resp []pcli.Flag
		Ω(rpc.GetFlags(new(interface{}), &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlFlags))
	})
})
