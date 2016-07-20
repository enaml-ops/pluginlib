package pcli

import (
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Given CLI flags", func() {
	Describe("given a new string slice flag", func() {
		Context("when there is a non empty stringslice value", func() {
			var ss Flag
			var controlString = "blah"
			BeforeEach(func() {
				ss = Flag{FlagType: StringSliceFlag}
				ss.Value = controlString
			})
			It("then it should set the default value string on the returned interface", func() {
				rval := ss.ToCli().(cli.StringSliceFlag)
				Ω(rval.Value.Value()).Should(ConsistOf(controlString))
			})
		})
		Context("when there is not a value", func() {
			var ss Flag
			BeforeEach(func() {
				ss = Flag{FlagType: StringSliceFlag}
			})
			It("then it should leave the Value empty on the returned type", func() {
				rval := ss.ToCli().(cli.StringSliceFlag)
				Ω(rval.Value.Value()).Should(BeEmpty())
			})
		})
	})
})
