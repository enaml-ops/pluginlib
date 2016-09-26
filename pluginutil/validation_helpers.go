package pluginutil

import (
	"github.com/xchapter7x/lo"
	"gopkg.in/urfave/cli.v2"
)

func CheckRequiredIntFlags(c *cli.Context, requiredFlags []string) []string {
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if c.Int(name) == 0 {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.String(name))
		}
	}
	return invalidNames
}

func CheckRequiredStringFlags(c *cli.Context, requiredFlags []string) []string {
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if c.String(name) == "" {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.String(name))
		}
	}
	return invalidNames
}

func CheckRequiredStringSliceFlags(c *cli.Context, requiredFlags []string) []string {
	invalidNames := []string{}
	for _, name := range requiredFlags {
		if len(c.StringSlice(name)) == 0 {
			invalidNames = append(invalidNames, name)
		} else {
			lo.G.Debug(name, "==>", c.StringSlice(name))
		}
	}
	return invalidNames
}
