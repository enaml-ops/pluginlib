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

func (dummyStore) Get(path, key string) (string, error)          { return "", nil }
func (dummyStore) Post(path, key, value string) error            { return nil }
func (dummyStore) PostBulk(path, values map[string]string) error { return nil }

func (dummyStore) GetBulk(path string) (map[string]string, error) {
	return map[string]string{
		"flag-one":  flagOne,
		"flag-two":  flagTwo,
		"flag-four": "asdf",
	}, nil
}

var _ = Describe("cred store", func() {
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
