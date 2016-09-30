package cred

import (
	"encoding/json"
	"fmt"
	"os"
)

type fileStore struct{}

// NewFileStore creates a Store backed by local files.
func NewFileStore() Store {
	return fileStore{}
}

// Get gets a single value from the specified path.
func (fs fileStore) Get(path, key string) (string, error) {
	kvPairs, err := readFile(path)
	if err != nil {
		return "", err
	}

	if val, ok := kvPairs[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("cred: %s not found in %s", key, path)
}

// GetBulk gets all key/value pairs from the specified path.
func (fs fileStore) GetBulk(path string) (map[string]string, error) {
	return readFile(path)
}

// Post updates a single value at the specified path.
func (fs fileStore) Post(path, key, value string) error {
	kvPairs, err := readFile(path)
	if err != nil {
		return err
	}
	kvPairs[key] = value
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(kvPairs)
}

// PostBulk updates all key/value pairs at the specified path.
func (fs fileStore) PostBulk(path string, values map[string]string) error {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(values)
}

func readFile(name string) (map[string]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	kvPairs := make(map[string]string)
	err = json.NewDecoder(f).Decode(&kvPairs)
	if err != nil {
		return nil, err
	}
	return kvPairs, nil
}
