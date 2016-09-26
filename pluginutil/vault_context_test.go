package pluginutil_test

import (
	"io/ioutil"
	"net/http"

	"github.com/enaml-ops/pluginlib/pcli"
	. "github.com/enaml-ops/pluginlib/pluginutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("given vault context", func() {
	Describe("given RotateSecrets", func() {
		Context("when not properly initialized and targetted or bad call", func() {

			var server *ghttp.Server
			var vault VaultRotater

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusBadRequest, b),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf")
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when called with a vault hash", func() {
				var err error

				BeforeEach(func() {
					err = vault.RotateSecrets("secret/move-along-secret", []byte(``))
				})

				It("then it should yield an error", func() {
					Ω(err).Should(HaveOccurred())
				})
			})
		})

		Context("when properly initialized and targetted", func() {
			var server *ghttp.Server
			var vault VaultRotater

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, b),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf")
			})

			AfterEach(func() {
				server.Close()
			})

			Context("when called with a vault hash", func() {
				var err error
				BeforeEach(func() {
					err = vault.RotateSecrets("secret/move-along-secret", []byte(``))
				})
				It("then it should populate the given vault hash with the given secrets", func() {
					Ω(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})

	Describe("given UnmarshalFlags", func() {

		Context("when unmarshalling string flags", func() {
			var server *ghttp.Server
			var vault *VaultUnmarshal

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vault.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, b),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf")
			})

			AfterEach(func() {
				server.Close()
			})

			It("should populate flags that were not specified on the command line", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx := NewContext([]string{"mycoolapp"}, ToCliFlagArray(flgs))
				Ω(ctx.String("knock")).Should(Equal("knocks"))
			})

			It("should not populate flags that aren't specified in the whitelist (UnmarshalSomeFlags)", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock", Value: "overwriteme"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "otherstuff", Value: "dontoverwriteme"},
				}
				vault.UnmarshalSomeFlags("secret/move-along-nothing-to-see-here", flgs, "knock")

				for i := range flgs {
					switch flgs[i].Name {
					case "knock":
						Ω(flgs[i].Value).Should(Equal("knocks"))
					case "otherstuff":
						Ω(flgs[i].Value).Should(Equal("dontoverwriteme"))
					}
				}
			})

			It("should not overwrite flags that were specified on the command line", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "knock"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx := NewContext([]string{"mycoolapp", "--knock", "something-different"}, ToCliFlagArray(flgs))
				Ω(ctx.String("knock")).Should(Equal("something-different"))
			})
		})

		Context("when unmarshalling string slice flags", func() {
			var server *ghttp.Server
			var vault *VaultUnmarshal

			BeforeEach(func() {
				b, _ := ioutil.ReadFile("fixtures/vaultslice.json")
				server = ghttp.NewServer()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, b),
					),
				)
				vault = NewVaultUnmarshal(server.URL(), "lkjaslkdjflkasjdf")
			})

			AfterEach(func() {
				server.Close()
			})

			It("should populate slice flags that were not specified on the command line", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "knock-slice"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "stuff"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx := NewContext([]string{"mycoolapp", "--stuff", "with-val"}, ToCliFlagArray(flgs))
				Ω(ctx.StringSlice("knock-slice")).Should(ConsistOf("knocks-1", "knocks-2", "knocks-3"))
			})

			It("should not populate slice flags that were specified on the command line", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringSliceFlag, Name: "knock-slice"},
					pcli.Flag{FlagType: pcli.StringFlag, Name: "stuff"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx := NewContext([]string{"mycoolapp", "--stuff", "with-val", "--knock-slice", "no-knocks-here"}, ToCliFlagArray(flgs))
				Ω(ctx.StringSlice("knock-slice")).Should(ConsistOf("no-knocks-here"))
			})

			It("should not populate flags that aren't defined in the context", func() {
				flgs := []pcli.Flag{
					pcli.Flag{FlagType: pcli.StringFlag, Name: "badda"},
				}
				vault.UnmarshalFlags("secret/move-along-nothing-to-see-here", flgs)
				ctx := NewContext([]string{"mycoolapp"}, ToCliFlagArray(flgs))
				Ω(ctx.String("knock")).Should(BeEmpty())
			})
		})
	})
})
