package cred

import (
	"errors"
	"fmt"
	"strings"

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

// NewStore creates a new Store based on the specified connection string.
// The following connection strings are supported:
//   - Hashicorp Vault: 'vault://TOKEN@domain:port'
//   - Filesystem: 'file://rootdir'
//
func NewStore(conn string) (Store, error) {
	store, details, err := parseConnString(conn)
	if err != nil {
		return nil, err
	}
	switch store {
	case "vault":
		tokenSep := strings.Index(details, "@")
		if tokenSep == -1 {
			return nil, fmt.Errorf("invalid Vault connection string: %q should be TOKEN@domain:port", details)
		}
		token := details[:tokenSep]
		domain := details[tokenSep+len("@"):]
		return NewVaultStore(domain, token), nil
	case "file":
		return NewFileStore(details), nil
	default:
		return nil, fmt.Errorf("unknown cred store %q", store)
	} 
}

func parseConnString(conn string) (store, details string, err error) {
	idx := strings.Index(conn, "://")
	if idx == -1 {
		return "", "", errors.New("invalid connection string, missing '://'")
	}
	store = conn[:idx]
	details = conn[idx+len("://"):]
	return store, details, nil
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
