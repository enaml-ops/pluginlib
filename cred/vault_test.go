package cred_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/enaml-ops/pluginlib/cred"
)

var _ = Describe("Vault credential store", func() {
	Context("get (happy path)", func() {
		var (
			server *ghttp.Server
			store  cred.Store
		)

		BeforeEach(func() {
			b, err := ioutil.ReadFile("fixtures/vault.json")
			Ω(err).ShouldNot(HaveOccurred())
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, b),
				),
			)
			store = cred.NewVaultStore(server.URL(), "token")
		})

		AfterEach(func() {
			server.Close()
		})

		It("can get all values at a path", func() {
			props, err := store.GetBulk("path")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(props).Should(HaveKeyWithValue("knock", "knocks"))
			Ω(props).Should(HaveKeyWithValue("otherstuff", "knock"))
		})

		It("can get a single value", func() {
			otherstuff, err := store.Get("path", "otherstuff")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(otherstuff).Should(Equal("knock"))
		})

		It("errors when getting a value that doesn't exist", func() {
			_, err := store.Get("path", "asdfasdfa")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("when vault returns malformed JSON", func() {
		var (
			server *ghttp.Server
			store  cred.Store
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, `{"lease_id"{:",renewa,ble":false,"lease_duration":2592000,"data":{"knock":"knocks",}`),
				),
			)
			store = cred.NewVaultStore(server.URL(), "token")
		})

		AfterEach(func() {
			server.Close()
		})

		It("returns an error", func() {
			_, err := store.GetBulk("path")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("post (happy path)", func() {
		const (
			expectedJSON = `{"knock": "knocks", "otherstuff": "knock"}`
			token        = "asdfasdfasdfasdfasdf"
		)
		var (
			server *ghttp.Server
			store  cred.Store
		)

		BeforeEach(func() {
			b, err := ioutil.ReadFile("fixtures/vault.json")
			Ω(err).ShouldNot(HaveOccurred())
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusOK, b),
					ghttp.VerifyRequest("POST", "/v1/myhash"), // verify that we got a POST
					ghttp.VerifyJSON(expectedJSON),            // verify the contents of the body
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyHeaderKV("X-Vault-Token", token),
				),
			)
			store = cred.NewVaultStore(server.URL(), token)
		})

		AfterEach(func() {
			server.Close()
		})

		It("posts the correct JSON data", func() {
			props := map[string]string{
				"knock":      "knocks",
				"otherstuff": "knock",
			}
			Ω(store.PostBulk("myhash", props)).Should(Succeed())
		})
	})

	Context("when vault return a bad HTTP status code", func() {
		var (
			server *ghttp.Server
			store  cred.Store
		)

		BeforeEach(func() {
			b, err := ioutil.ReadFile("fixtures/vault.json")
			Ω(err).ShouldNot(HaveOccurred())
			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusNotFound, b),
				),
			)
			store = cred.NewVaultStore(server.URL(), "token")
		})

		AfterEach(func() {
			server.Close()
		})

		It("returns an error", func() {
			values := map[string]string{
				"foo": "1",
				"bar": "2",
			}
			Ω(store.PostBulk("path", values)).ShouldNot(Succeed())
		})
	})
})
