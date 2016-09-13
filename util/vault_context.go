package pluginutil

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/xchapter7x/lo"
)

type VaultUnmarshaler interface {
	UnmarshalFlags(hash string, flgs []pcli.Flag) (err error)
}

type VaultRotater interface {
	RotateSecrets(hash string, secrets interface{}) error
}

func NewVaultUnmarshal(domain, token string) *VaultUnmarshal {
	return &VaultUnmarshal{
		Domain: domain,
		Token:  token,
		client: defaultClient(),
	}
}

type VaultUnmarshal struct {
	Domain string
	Token  string
	client *http.Client
}

type vaultJsonObject struct {
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration float64           `json:"lease_duration"`
	Data          map[string]string `json:"data"`
	Warnings      interface{}       `json:"warnings"`
	Auth          interface{}       `json:"auth"`
}

func (s *VaultUnmarshal) RotateSecrets(hash string, secrets interface{}) error {
	return s.setVaultHashValues(hash, secrets.([]byte))
}

func (s *VaultUnmarshal) UnmarshalFlags(hash string, flgs []pcli.Flag) error {
	b := s.getVaultHashValues(hash)
	vaultObj := new(vaultJsonObject)
	if err := json.Unmarshal(b, vaultObj); err != nil {
		return err
	}

	for i := range flgs {
		flagName := flgs[i].Name
		if vaultValue, ok := vaultObj.Data[flagName]; ok {
			flgs[i].Value = vaultValue
			lo.G.Debugf("set %s flag from vault (value=%s)", flagName, vaultValue)
		}
	}
	return nil
}

func (s *VaultUnmarshal) setVaultHashValues(hash string, body []byte) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/%s", s.Domain, hash), bytes.NewBuffer(body))
	if err != nil {
		lo.G.Errorf("error in vault request %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	s.decorateWithToken(req)

	res, err := s.client.Do(req)
	if err != nil {
		lo.G.Errorf("error calling client %v", err)
		return err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		lo.G.Errorf("bad resp code from vault: %d", res.StatusCode)
		return fmt.Errorf("status code is not ok: %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		lo.G.Errorf("error reading response: %v", err)
		return err
	}

	lo.G.Debugf("vault res: %v", string(b))
	return nil
}

func (s *VaultUnmarshal) getVaultHashValues(hash string) []byte {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s", s.Domain, hash), nil)
	s.decorateWithToken(req)
	res, _ := s.client.Do(req)
	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func (s *VaultUnmarshal) decorateWithToken(req *http.Request) {
	req.Header.Add("X-Vault-Token", s.Token)
}

func defaultClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}
