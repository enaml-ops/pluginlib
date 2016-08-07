package pluginutil_test

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/enaml-ops/pluginlib/pcli"
	. "github.com/enaml-ops/pluginlib/util"
	"github.com/enaml-ops/pluginlib/util/utilfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("given: a VaultUnmarshal", func() {

	Context("when not properly initialized and targetted or bad call", func() {
		Describe("given rotate secrets", func() {
			var server *ghttp.Server
			var vault VaultRotater

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusBadRequest, string(b)),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf", DefaultClient())
			})

			AfterEach(func() {
				server.Close()
			})
			Context("when called with a vault hash", func() {
				var err error
				BeforeEach(func() {
					err = vault.RotateSecrets("secret/move-along-secret", nil)
				})
				It("then it should yield an error", func() {
					Ω(err).Should(HaveOccurred())
				})
			})
		})
	})
	Context("when properly initialized and targetted", func() {
		Describe("given rotate secrets", func() {
			var server *ghttp.Server
			var vault VaultRotater

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, string(b)),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf", DefaultClient())
			})

			AfterEach(func() {
				server.Close()
			})
			Context("when called with a vault hash", func() {
				var err error
				BeforeEach(func() {
					err = vault.RotateSecrets("secret/move-along-secret", nil)
				})
				It("then it should populate the given vault hash with the given secrets", func() {
					Ω(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("given a defaultclient", func() {
		var server *ghttp.Server
		var vault VaultUnmarshaler

		BeforeEach(func() {
			b, _ := ioutil.ReadFile("fixtures/vault.json")
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, string(b)),
				),
			)
			vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf", DefaultClient())
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when calling unmarshalflags on a context that was not given the flag value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp"}, ToCliFlagArray(flgs))
			})

			It("should set the value in the flag using the given vault hash", func() {
				Ω(ctx.String("knock")).Should(Equal("knocks"))
			})
		})
	})

	Describe("given a properly initialized vaultoverlay targeting a vault", func() {
		var vault VaultUnmarshaler

		BeforeEach(func() {
			doer := new(utilfakes.FakeDoer)
			b, _ := os.Open("fixtures/vault.json")
			doer.DoReturns(&http.Response{
				Body: b,
			}, nil)
			vault = NewVaultUnmarshal("domain.com", "my-really-long-token", doer)
		})

		Context("when calling unmarshalflags on a context that was not given the flag value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp"}, ToCliFlagArray(flgs))
			})

			It("should set the value in the flag using the given vault hash", func() {
				Ω(ctx.String("knock")).Should(Equal("knocks"))
			})
		})

		Context("when calling unmarshalflags on a context that was given the flag value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp", "--knock", "something-different"}, ToCliFlagArray(flgs))
			})

			It("should overwrite the default vault value with the cli flag value given", func() {
				Ω(ctx.String("knock")).ShouldNot(Equal("knocks"))
				Ω(ctx.String("knock")).Should(Equal("something-different"))
			})
		})

		Context("when calling unmarshalflags on a context that was given a stringslice value from the cli", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				doer := new(utilfakes.FakeDoer)
				b, _ := os.Open("fixtures/vaultslice.json")

				doer.DoReturns(&http.Response{
					Body: b,
				}, nil)
				vault = NewVaultUnmarshal("domain.com", "my-really-long-token", doer)
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "knock-slice"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "stuff"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp", "--stuff", "with-val"}, ToCliFlagArray(flgs))
			})

			It("should overwrite the value in the flag using the given vault hash", func() {
				Ω(ctx.StringSlice("knock-slice")).Should(ConsistOf("knocks-1", "knocks-2", "knocks-3"))
			})
		})

		Context("when calling unmarshalflags on a context which was not defined with the flag contained in vault", func() {
			var ctx *cli.Context

			BeforeEach(func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "badda"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx = NewContext([]string{"mycoolapp"}, ToCliFlagArray(flgs))
			})

			It("then it should not set or create the flag in the context", func() {
				Ω(ctx.String("knock")).Should(BeEmpty())
			})
		})
	})
})
