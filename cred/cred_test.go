package cred_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/enaml-ops/pluginlib/cred"
	"github.com/enaml-ops/pluginlib/pcli"
)

const (
	flagOne   = "one"
	flagTwo   = "two"
	flagThree = "three"
)

type dummyStore struct{}

func (dummyStore) Get(path, key string) (string, error)                 { return "", nil }
func (dummyStore) Post(path, key, value string) error                   { return nil }
func (dummyStore) PostBulk(path string, values map[string]string) error { return nil }

func (dummyStore) GetBulk(path string) (map[string]string, error) {
	return map[string]string{
		"flag-one":  flagOne,
		"flag-two":  flagTwo,
		"flag-four": "asdf",
	}, nil
}

var _ = Describe("cred store", func() {
	Context("NewStore", func() {
		It("returns an error when given an unsupported connection string", func() {
			_, err := cred.NewStore("unsupported://connection-string")
			Ω(err).Should(HaveOccurred())

			_, err = cred.NewStore("not-even-a-connection-string")
			Ω(err).Should(HaveOccurred())
		})

		It("returns an error when given an invalid Vault connection string", func() {
			_, err := cred.NewStore("vault://myvaultwithoutatoken.com:8200")
			Ω(err).Should(HaveOccurred())

			_, err = cred.NewStore("vault://VAULT_TOKEN_BUT_NO_DOMAIN")
			Ω(err).Should(HaveOccurred())
		})

		It("creates a Vault store when given a valid connection string", func() {
			store, err := cred.NewStore("vault://token@10.0.1.2:8200")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(store).ShouldNot(BeNil())
			// TODO: move into package cred and type assert the result?
		})

		It("creates a filesystem-backed store", func() {
			store, err := cred.NewStore("file://.")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(store).ShouldNot(BeNil())
			// TODO: move into package cred and type assert the result?
		})
	})

	Context("Overlay", func() {
		var (
			flags []pcli.Flag
		)

		BeforeEach(func() {
			flags = []pcli.Flag{
				pcli.CreateStringFlag("flag-one", "", "test"),
				pcli.CreateStringFlag("flag-two", ""),
				pcli.CreateStringFlag("flag-three", "", flagThree),
			}
			cred.Overlay("", flags, dummyStore{})
		})

		It("sets the values of matching flags", func() {
			for _, flag := range flags {
				switch flag.Name {
				case "flag-one":
					Ω(flag.Value).Should(Equal(flagOne))
				case "flag-two":
					Ω(flag.Value).Should(Equal(flagTwo))
				}
			}
		})

		It("doesn't modify non-matching flags", func() {
			for _, flag := range flags {
				if flag.Name == "flag-three" {
					Ω(flag.Value).Should(Equal(flagThree))
				}
			}
		})
	})
})
