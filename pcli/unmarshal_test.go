package pcli_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/util"

	cli "gopkg.in/urfave/cli.v2"
)

var _ = Describe("Unmarshal flags", func() {

	var (
		flags   []cli.Flag
		context *cli.Context
	)

	Context("when the context has all required flags", func() {

		BeforeEach(func() {
			flags = []cli.Flag{
				&cli.StringFlag{Name: "string-flag"},
				&cli.IntFlag{Name: "int-flag"},
				&cli.StringSliceFlag{Name: "string-slice-flag"},
				&cli.IntSliceFlag{Name: "int-slice-flag"},
				&cli.BoolFlag{Name: "bool-flag"},
				&cli.Float64Flag{Name: "float-flag"},
				&cli.StringFlag{Name: "dont-change"},
				&cli.StringFlag{Name: "dummy"},
				&cli.StringFlag{Name: "dummy2"},
			}
			context = pluginutil.NewContext([]string{
				"foo",
				"--string-flag", "mystring",
				"--int-flag", "42",
				"--string-slice-flag", "value1", "--string-slice-flag", "value2",
				"--int-slice-flag", "1", "--int-slice-flag", "2", "--int-slice-flag", "3",
				"--bool-flag",
				"--float-flag", "1.23",
				"--dont-change", "me",
			}, flags)
		})

		Context("when unmarshalling into a non-struct type", func() {
			It("should panic", func() {
				f := func() {
					var x string
					pcli.UnmarshalFlags(x, context)
				}
				Ω(f).Should(Panic())
			})
		})

		Context("when marshalling into a value and not a pointer", func() {
			type FlagTest struct {
				StringFlag      string
				IntFlag         int
				StringSliceFlag []string
				IntSliceFlag    []int
				BoolFlag        bool
				FloatFlag       float64
			}

			It("should panic", func() {
				f := func() {
					var t FlagTest
					pcli.UnmarshalFlags(t, context)
				}
				Ω(f).Should(Panic())
			})
		})

		Context("when unmarshalling into a nil value", func() {
			It("should panic", func() {
				f := func() {
					pcli.UnmarshalFlags(nil, context)
				}
				Ω(f).Should(Panic())
			})
		})

		Context("when unmarshalling using default flag names", func() {
			type FlagTest struct {
				StringFlag      string
				IntFlag         int
				StringSliceFlag []string
				IntSliceFlag    []int
				BoolFlag        bool
				FloatFlag       float64
				DontChange      string `omg:"-"`
				Dummy1          string `omg:"-"`
				Dummy2          string
			}

			var t FlagTest

			BeforeEach(func() {
				t = FlagTest{}
				t.DontChange = "original"
				t.Dummy1 = "dummy"
				t.Dummy2 = "dummy"
				pcli.UnmarshalFlags(&t, context)
			})

			It("should populate the struct fields", func() {
				Ω(t.StringFlag).Should(Equal("mystring"))
				Ω(t.IntFlag).Should(Equal(42))
				Ω(t.StringSliceFlag).Should(ConsistOf("value1", "value2"))
				Ω(t.IntSliceFlag).Should(ConsistOf(1, 2, 3))
				Ω(t.BoolFlag).Should(BeTrue())
				Ω(t.FloatFlag).Should(Equal(1.23))
			})

			It("doesn't modify fields annotated with '-'", func() {
				Ω(t.DontChange).Should(Equal("original"))
				Ω(t.DontChange).ShouldNot(Equal("me"))
				Ω(t.Dummy1).Should(Equal("dummy"))
			})
		})

		Context("when unmarshalling embedded fields", func() {
			type (
				Slices struct {
					StringSliceFlag []string
					IntSliceFlag    []int
				}
				FlagTest struct {
					Slices
					StringFlag string
					IntFlag    int
					BoolFlag   bool
					FloatFlag  float64
				}
			)

			var t FlagTest

			BeforeEach(func() {
				t = FlagTest{}
				pcli.UnmarshalFlags(&t, context)
			})

			It("sets the non-embedded fields correctly", func() {
				Ω(t.StringFlag).Should(Equal("mystring"))
				Ω(t.IntFlag).Should(Equal(42))
				Ω(t.BoolFlag).Should(BeTrue())
				Ω(t.FloatFlag).Should(Equal(1.23))
			})

			It("sets the embedded fields correctly", func() {
				Ω(t.StringSliceFlag).Should(ConsistOf("value1", "value2"))
				Ω(t.IntSliceFlag).Should(ConsistOf(1, 2, 3))
			})
		})

		Context("when unmarshalling using struct tags", func() {
			type FlagTest struct {
				MyString      string   `omg:"string-flag"`
				MyInt         int      `omg:"int-flag"`
				MyStringSlice []string `omg:"string-slice-flag"`
				MyIntSlice    []int    `omg:"int-slice-flag"`
				MyBool        bool     `omg:"bool-flag"`
				MyFloat       float64  `omg:"float-flag"`
			}
			var t FlagTest

			BeforeEach(func() {
				t = FlagTest{}
				pcli.UnmarshalFlags(&t, context)
			})

			It("should populate all struct fields", func() {
				Ω(t.MyString).Should(Equal("mystring"))
				Ω(t.MyInt).Should(Equal(42))
				Ω(t.MyStringSlice).Should(ConsistOf("value1", "value2"))
				Ω(t.MyIntSlice).Should(ConsistOf(1, 2, 3))
				Ω(t.MyBool).Should(BeTrue())
				Ω(t.MyFloat).Should(Equal(1.23))
			})
		})
	})

	Context("when the context is missing required flags", func() {
		BeforeEach(func() {
			flags = []cli.Flag{
				&cli.StringFlag{Name: "string-flag"},
				&cli.IntFlag{Name: "int-flag"},
				&cli.StringSliceFlag{Name: "string-slice-flag"},
				&cli.IntSliceFlag{Name: "int-slice-flag"},
				&cli.BoolFlag{Name: "bool-flag"},
				&cli.Float64Flag{Name: "float-flag"},
			}
			context = pluginutil.NewContext([]string{
				"foo",
				"--string-slice-flag", "value1", "--string-slice-flag", "value2",
				"--int-slice-flag", "1", "--int-slice-flag", "2", "--int-slice-flag", "3",
				"--bool-flag",
				"--float-flag", "1.23",
			}, flags)
		})

		It("returns an error", func() {
			type FlagTest struct {
				StringFlag      string
				IntFlag         int
				StringSliceFlag []string
				IntSliceFlag    []int
				BoolFlag        bool
				FloatFlag       float64
			}
			t := FlagTest{}
			Ω(pcli.UnmarshalFlags(&t, context)).ShouldNot(Succeed())
		})

		Context("when the context is missing required slice flags", func() {
			BeforeEach(func() {
				flags = []cli.Flag{
					&cli.StringFlag{Name: "string-flag"},
					&cli.IntFlag{Name: "int-flag"},
					&cli.StringSliceFlag{Name: "string-slice-flag"},
					&cli.IntSliceFlag{Name: "int-slice-flag"},
					&cli.BoolFlag{Name: "bool-flag"},
					&cli.Float64Flag{Name: "float-flag"},
				}
				context = pluginutil.NewContext([]string{
					"foo",
					"--string-flag", "foo",
					"--int-flag", "42",
					"--int-slice-flag", "1", "--int-slice-flag", "2", "--int-slice-flag", "3",
					"--bool-flag",
					"--float-flag", "1.23",
				}, flags)
			})

			type FlagTest struct {
				StringFlag      string
				IntFlag         int
				StringSliceFlag []string
				IntSliceFlag    []int
				BoolFlag        bool
				FloatFlag       float64
			}

			It("should not panic", func() {
				Ω(func() {
					t := FlagTest{}
					Ω(pcli.UnmarshalFlags(&t, context)).ShouldNot(Succeed())
				}).ShouldNot(Panic())
			})

			It("should return an error", func() {
				t := FlagTest{}
				Ω(pcli.UnmarshalFlags(&t, context)).ShouldNot(Succeed())
			})
		})

		Context("when using flags with defaults", func() {
			BeforeEach(func() {
				flags = []cli.Flag{
					&cli.StringFlag{Name: "string-flag", Value: "foo"},
					&cli.IntFlag{Name: "int-flag", Value: 42},
				}
				context = pluginutil.NewContext([]string{"foo", "--string-flag", "bar"}, flags)
			})

			It("should not require the flags to be provided", func() {
				type FlagTest struct {
					StringFlag string
					IntFlag    int
				}
				t := FlagTest{}
				err := pcli.UnmarshalFlags(&t, context)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(t.StringFlag).Should(Equal("bar"))
				Ω(t.IntFlag).Should(Equal(42))
			})
		})

		Context("when unmarshalling into a struct with OPTIONAL flags", func() {
			BeforeEach(func() {
				flags = []cli.Flag{
					&cli.StringFlag{Name: "string-flag"},
					&cli.IntFlag{Name: "int-flag"},
				}
				context = pluginutil.NewContext([]string{"foo"}, flags)
			})

			It("doesn't issue an error if optional flags are missing", func() {
				type FlagTest struct {
					Str     string `omg:"string-flag,optional"`
					IntFlag int    `omg:"int-flag,optional"`
				}
				t := FlagTest{}
				Ω(pcli.UnmarshalFlags(&t, context)).Should(Succeed())
			})
		})

	})
})
