package pcli

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/urfave/cli.v2"
)

var _ = Describe("Given CLI flags", func() {
	Describe("given a Flag", func() {
		var flg Flag
		var flagname = "hi-there-flag"

		Context("when: calling tocli and there is not a envvar value set", func() {
			BeforeEach(func() {
				flg = Flag{FlagType: StringFlag, Name: flagname}
			})
			It("then it should automatically set the env var", func() {
				Ω(flg.ToCli().(*cli.StringFlag).EnvVars).ShouldNot(BeEmpty())
				Ω(flg.ToCli().(*cli.StringFlag).EnvVars).Should(ConsistOf("OMG_HI_THERE_FLAG"))
			})
		})

		Context("when: calling tocli and there is a envvar value already set", func() {
			var control = "bleh"
			BeforeEach(func() {
				flg = Flag{FlagType: StringFlag, Name: flagname, EnvVar: control}
			})
			It("then it should use that", func() {
				Ω(flg.ToCli().(*cli.StringFlag).EnvVars).ShouldNot(BeEmpty())
				Ω(flg.ToCli().(*cli.StringFlag).EnvVars).Should(ConsistOf(control))
			})
		})
	})

	Describe("given a new string slice flag", func() {
		Context("when there is a non empty stringslice value", func() {
			var ss Flag
			var controlString = "blah"
			BeforeEach(func() {
				ss = Flag{FlagType: StringSliceFlag}
				ss.Value = controlString
			})
			It("then it should set the default value string on the returned interface", func() {
				rval := ss.ToCli().(*cli.StringSliceFlag)
				Ω(rval.Value.Value()).Should(ConsistOf(controlString))
			})
		})
		Context("when there is not a value", func() {
			var ss Flag
			BeforeEach(func() {
				ss = Flag{FlagType: StringSliceFlag}
			})
			It("then it should leave the Value empty on the returned type", func() {
				rval := ss.ToCli().(*cli.StringSliceFlag)
				Ω(rval.Value.Value()).Should(BeEmpty())
			})
		})
	})
})
