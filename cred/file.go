package cred

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type fileStore struct {
	rootDir string
}

// NewFileStore creates a Store backed by local files at the specified root directory.
func NewFileStore(root string) Store {
	return fileStore{
		rootDir: root,
	}
}

// Get gets a single value from the specified path.
func (fs fileStore) Get(path, key string) (string, error) {
	path = filepath.Join(fs.rootDir, path)
	kvPairs, err := fs.readFile(path)
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
	path = filepath.Join(fs.rootDir, path)
	return fs.readFile(path)
}

// Post updates a single value at the specified path.
func (fs fileStore) Post(path, key, value string) error {
	path = filepath.Join(fs.rootDir, path)
	kvPairs, err := fs.readFile(path)
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
	path = filepath.Join(fs.rootDir, path)
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(values)
}

func (fs fileStore) readFile(name string) (map[string]string, error) {
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
