package registry_test

import (
	"runtime"
	"testing"

	. "github.com/enaml-ops/pluginlib/registry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Registry", func() {

	Describe("given RegisterProduct function", func() {
		Context("when called w/ valid parameters", func() {

			BeforeEach(func() {
				if testing.Short() {
					Skip("plugin registry tests skipped in short mode")
				}
				RegisterProduct("./fixtures/product/testproductplugin-" + runtime.GOOS)
			})

			It("then it should register the plugin from the given path in the registry", func() {
				products := ListProducts()
				Ω(len(products)).Should(Equal(1))
				Ω(products["myfakeproduct"]).ShouldNot(BeNil())
			})
		})
	})
	Describe("given RegisterCloudConfig function", func() {
		Context("when called w/ valid parameters", func() {

			BeforeEach(func() {
				if testing.Short() {
					Skip("plugin registry tests skipped in short mode")
				}
				RegisterCloudConfig("./fixtures/cloudconfig/testplugin-" + runtime.GOOS)
			})

			It("then it should register the plugin from the given path in the registry", func() {
				cloudconfigs := ListCloudConfigs()
				Ω(len(cloudconfigs)).Should(Equal(1))
				Ω(cloudconfigs["myfakecloudconfig"]).ShouldNot(BeNil())
			})
		})
	})
})
