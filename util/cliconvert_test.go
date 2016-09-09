package pluginutil_test

import (
	"gopkg.in/urfave/cli.v2"
	"github.com/enaml-ops/pluginlib/pcli"
	. "github.com/enaml-ops/pluginlib/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("given ToCliFlagArray", func() {
	Context("when called with a []pcli.Flag", func() {
		controlFlags := []pcli.Flag{
			pcli.Flag{FlagType: pcli.StringFlag, Name: "blahstring"},
			pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "blahslice"},
			pcli.Flag{FlagType: pcli.IntFlag, Name: "blahint"},
			pcli.Flag{FlagType: pcli.BoolFlag, Name: "blahbool"},
			pcli.Flag{FlagType: pcli.BoolTFlag, Name: "blahboolt"},
		}
		It("then it should convert to a []cli.Flag", func() {
			cliFlags := ToCliFlagArray(controlFlags)
			Ω(cliFlags).ShouldNot(Equal(make([]cli.Flag, 0)))
			Ω(len(cliFlags)).Should(Equal(len(controlFlags)))
		})
	})
})
