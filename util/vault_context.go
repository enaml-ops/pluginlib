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

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

func NewVaultUnmarshal(domain, token string, client Doer) *VaultUnmarshal {
	return &VaultUnmarshal{
		Domain: domain,
		Token:  token,
		Client: client,
	}
}

type VaultUnmarshal struct {
	Domain string
	Token  string
	Client Doer
}

type VaultJsonObject struct {
	LeaseID       string            `json:"lease_id"`
	Renewable     bool              `json:"renewable"`
	LeaseDuration float64           `json:"lease_duration"`
	Data          map[string]string `json:"data"`
	Warnings      interface{}       `json:"warnings"`
	Auth          interface{}       `json:"auth"`
}

func (s *VaultUnmarshal) RotateSecrets(hash string, secrets interface{}) (err error) {
	return s.setVaultHashValues(hash, secrets.([]byte))
}

func (s *VaultUnmarshal) UnmarshalFlags(hash string, flgs []pcli.Flag) (err error) {
	b := s.getVaultHashValues(hash)
	vaultObj := new(VaultJsonObject)
	json.Unmarshal(b, vaultObj)

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
	var err error
	var req *http.Request
	var res *http.Response
	var b []byte

	if req, err = http.NewRequest("POST", fmt.Sprintf("%s/v1/%s", s.Domain, hash), bytes.NewBuffer(body)); err != nil {
		lo.G.Errorf("error in vault request %v", err)

	} else {
		req.Header.Set("Content-Type", "application/json")
		s.decorateWithToken(req)

		if res, err = s.Client.Do(req); err != nil {
			lo.G.Errorf("error calling client %v", err)

		} else if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
			err = fmt.Errorf("status code is not ok: %v", res.StatusCode)
			lo.G.Error(err.Error())

		} else {

			if b, err = ioutil.ReadAll(res.Body); err != nil {
				lo.G.Errorf("error in reading response %v", err)
			}
			lo.G.Debugf("vault res: %v", string(b))
		}
	}
	return err
}

func (s *VaultUnmarshal) getVaultHashValues(hash string) []byte {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/v1/%s", s.Domain, hash), nil)
	s.decorateWithToken(req)
	res, _ := s.Client.Do(req)
	b, _ := ioutil.ReadAll(res.Body)
	return b
}

func (s *VaultUnmarshal) decorateWithToken(req *http.Request) {
	req.Header.Add("X-Vault-Token", s.Token)
}

func DefaultClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}
