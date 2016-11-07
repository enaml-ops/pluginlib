package cred_test

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/enaml-ops/pluginlib/cred"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("file-backed cred store", func() {
	var cs cred.Store

	BeforeEach(func() {
		cs = cred.NewFileStore(".")
	})

	Context("when targetting a non-existent path", func() {
		It("returns an error", func() {
			_, err := cs.Get("this/path/does/not/exist", "foo")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("reading single values", func() {
		It("gets the correct values", func() {
			secret1, err := cs.Get("fixtures/file.json", "pass1")
			Ω(err).ShouldNot(HaveOccurred())

			secret2, err := cs.Get("fixtures/file.json", "pass2")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(secret1).Should(Equal("secret1"))
			Ω(secret2).Should(Equal("secret2"))
		})

		It("returns an error if a value is not found", func() {
			val, err := cs.Get("fixtures/file.json", "foo")
			Ω(err).Should(HaveOccurred())
			Ω(val).Should(Equal(""))
		})
	})

	Context("reading multiple values", func() {
		It("gets all values", func() {
			val, err := cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(val).Should(HaveKeyWithValue("pass1", "secret1"))
			Ω(val).Should(HaveKeyWithValue("pass2", "secret2"))
		})
	})

	Context("writing a single value", func() {
		var (
			tmp    *os.File
			backup string
		)

		BeforeEach(func() {
			var err error
			tmp, err = ioutil.TempFile(".", "pluginlib-file-test")
			Ω(err).ShouldNot(HaveOccurred())
			backup = tmp.Name()
			defer tmp.Close()

			orig, err := os.Open("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			defer orig.Close()

			io.Copy(tmp, orig)
		})

		AfterEach(func() {
			err := os.Rename(backup, "fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("returns an error when targetting a bad path", func() {
			Ω(cs.Post("this/path/does/not/exist", "foo", "bar")).ShouldNot(Succeed())
		})

		It("overwrites existing values", func() {
			orig, err := cs.Get("fixtures/file.json", "pass1")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(orig).Should(Equal("secret1"))

			err = cs.Post("fixtures/file.json", "pass1", "newsecret1")
			Ω(err).ShouldNot(HaveOccurred())

			pass, err := cs.Get("fixtures/file.json", "pass1")
			Ω(err).ShouldNot(HaveOccurred())

			Ω(pass).Should(Equal("newsecret1"))
		})

		It("adds new values", func() {
			vals, err := cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vals).Should(HaveLen(2))

			Ω(cs.Post("fixtures/file.json", "pass3", "secret3")).Should(Succeed())

			vals, err = cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vals).Should(HaveLen(3))
			Ω(vals).Should(HaveKeyWithValue("pass3", "secret3"))
		})
	})

	Context("writing multiple values", func() {
		var (
			tmp    *os.File
			backup string
		)

		BeforeEach(func() {
			var err error
			tmp, err = ioutil.TempFile(".", "pluginlib-file-test")
			Ω(err).ShouldNot(HaveOccurred())
			backup = tmp.Name()
			defer tmp.Close()

			orig, err := os.Open("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			defer orig.Close()

			io.Copy(tmp, orig)
		})

		AfterEach(func() {
			err := os.Rename(backup, "fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("updates existing values", func() {
			newVals := map[string]string{
				"pass1": "newsecret1",
				"pass2": "newsecret2",
			}
			Ω(cs.PostBulk("fixtures/file.json", newVals)).Should(Succeed())

			vals, err := cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vals).Should(HaveLen(2))
			Ω(vals).Should(HaveKeyWithValue("pass1", "newsecret1"))
			Ω(vals).Should(HaveKeyWithValue("pass2", "newsecret2"))
		})

		It("overwrites all values", func() {
			vals, err := cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vals).Should(HaveLen(2))

			newVals := map[string]string{
				"pass1": "newsecret1",
			}
			Ω(cs.PostBulk("fixtures/file.json", newVals)).Should(Succeed())

			vals, err = cs.GetBulk("fixtures/file.json")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(vals).Should(HaveLen(1))
			Ω(vals).Should(HaveKeyWithValue("pass1", "newsecret1"))
		})
	})
})
