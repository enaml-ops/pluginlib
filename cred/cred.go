package cred

import (
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
)

// Store is a repository of credentials for use by omg plugins.
type Store interface {
	// Get gets a single value from the specified path.
	Get(path, key string) (string, error)

	// GetBulk gets all key/value pairs from the specified path.
	GetBulk(path string) (map[string]string, error)

	// Post updates a single value at the specified path.
	Post(path, key, value string) error

	// PostBulk updates all key/value pairs at the specified path.
	PostBulk(path string, values map[string]string) error
}

// Overlay provides default values for the specified flags
// using matching values from a credential store.
func Overlay(path string, flags []pcli.Flag, store Store) error {
	props, err := store.GetBulk(path)
	if err != nil {
		return err
	}
	for i := range flags {
		name := flags[i].Name
		if val, ok := props[name]; ok {
			flags[i].Value = val
			lo.G.Debugf("set %s flag from cred store (value=%s)", name, val)
		}
	}
	return nil
}
