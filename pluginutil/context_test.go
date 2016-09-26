package pluginutil_test

import (
	. "github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("NewContext function", func() {
	Context("when called with valid args and flags", func() {
		It("then it should return a properly init'd cli.context", func() {
			ctx := NewContext([]string{"test", "--this", "that"}, []cli.Flag{
				&cli.StringFlag{Name: "this"},
			})
			Ω(ctx.String("this")).Should(Equal("that"))
		})
	})
})
