package cred

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type vaultStore struct {
	domain string
	token  string
}

// NewVaultStore creates a Store backed by Hashicorp's Vault.
func NewVaultStore(domain, token string) Store {
	return &vaultStore{
		domain: domain,
		token:  token,
	}
}

// Get gets a single value from the specified path.
func (v *vaultStore) Get(path, key string) (string, error) {
	props, err := v.GetBulk(path)
	if err != nil {
		return "", err
	}
	if val, ok := props[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("cred: %s not found in %s", key, path)
}

// GetBulk gets all key/value pairs from the specified path.
func (v *vaultStore) GetBulk(path string) (map[string]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s", v.domain, path), nil)
	if err != nil {
		return nil, err
	}
	v.decorateWithToken(req)
	client := v.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var js vaultJSON
	err = json.NewDecoder(resp.Body).Decode(&js)
	if err != nil {
		return nil, err
	}
	return js.Data, nil
}

// Post updates a single value at the specified path.
func (v *vaultStore) Post(path, key, value string) error {
	// vault doesn't support writing a single value,
	// so we do a read-modify-write operation
	props, err := v.GetBulk(path)
	if err != nil {
		return err
	}

	props[key] = value
	return v.PostBulk(path, props)
}

// PostBulk updates all key/value pairs at the specified path.
func (v *vaultStore) PostBulk(path string, values map[string]string) error {
	body, err := json.Marshal(values)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/%s", v.domain, path), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	v.decorateWithToken(req)

	client := v.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// we aren't concerned with the body for now
	io.Copy(ioutil.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("cred: vault POST failed with status %d", resp.StatusCode)
	}

	return nil
}

type vaultJSON struct {
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration float64           `json:"lease_duration"`
	Data          map[string]string `json:"data"`
	Warnings      interface{}       `json:"warnings"`
	Auth          interface{}       `json:"auth"`
}

func (v *vaultStore) getClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func (v *vaultStore) decorateWithToken(req *http.Request) {
	req.Header.Add("X-Vault-Token", v.token)
}
