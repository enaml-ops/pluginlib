package pluginutil_test

import (
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/enaml-ops/pluginlib/util"
)

var _ = Describe("CloudConfigInfer", func() {
	var inferer *CloudConfigInfer
	Describe("given it is properly initialized with a valid and complete cloudconfig manifest", func() {
		BeforeEach(func() {
			b, _ := ioutil.ReadFile("fixtures/cloudconfig.yml")
			inferer = NewCloudConfigInferFromBytes(b)
			inferer.InferDefaultVMType()
		})
		testInfer("InferDefaultVMType", "small", "then it should grab the first it finds from cld cnf", func() string { return inferer.InferDefaultVMType() })
		testInfer("InferDefaultDiskType", "small", "then it should grab the first it finds from cld cnf", func() string { return inferer.InferDefaultDiskType() })
		testInfer("InferDefaultNetwork", "private", "then it should grab the first it finds from cld cnf", func() string { return inferer.InferDefaultNetwork() })
		testInfer("InferDefaultAZ", "z1,z2", "then it should grab all available as a csv", func() string { return inferer.InferDefaultAZ() })
	})
})

func testInfer(methodname, control, bahave string, method func() string) {
	Context(fmt.Sprintf("when the %s is called", methodname), func() {
		var vmtypename string

		BeforeEach(func() {
			vmtypename = method()
		})

		It(bahave, func() {
			Ω(vmtypename).Should(Equal(control))
		})
	})
}
