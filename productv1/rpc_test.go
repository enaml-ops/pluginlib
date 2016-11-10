package product_test

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/productv1"
	"github.com/enaml-ops/pluginlib/productv1/productv1fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("productv1 RPC", func() {
	var d *productv1fakes.FakeDeployer

	BeforeEach(func() {
		d = new(productv1fakes.FakeDeployer)
	})

	It("Forwards calls to GetMeta", func() {
		controlMeta := product.Meta{
			Name: "fakemeta",
		}
		d.GetMetaReturns(controlMeta)
		rpc := product.RPCServer{
			Impl: d,
		}

		var resp product.Meta
		Ω(rpc.GetMeta(new(interface{}), &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlMeta))
	})

	It("Forwards calls to GetProduct", func() {
		controlResp := product.Response{
			Bytes:  []byte{0, 1, 2},
			ErrRes: "",
		}
		d.GetProductReturns(controlResp.Bytes, nil)
		rpc := product.RPCServer{Impl: d}

		var resp product.Response
		Ω(rpc.GetProduct(product.Args{
			Args:        []string{"product"},
			CloudConfig: []byte{0, 1, 2},
			CredStore:   nil,
		}, &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlResp))
	})

	It("Forwards calls to GetFlags", func() {
		controlFlags := []pcli.Flag{
			pcli.CreateStringFlag("str", "dummy", ""),
		}
		d.GetFlagsReturns(controlFlags)
		rpc := product.RPCServer{
			Impl: d,
		}

		var resp []pcli.Flag
		Ω(rpc.GetFlags(new(interface{}), &resp)).Should(Succeed())
		Ω(resp).Should(Equal(controlFlags))
	})
})
