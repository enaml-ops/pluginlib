package registry

import (
	"log"
	"os/exec"

	"github.com/enaml-ops/pluginlib/cloudconfig"
	"github.com/enaml-ops/pluginlib/pcli"
	"github.com/enaml-ops/pluginlib/product"
	"github.com/hashicorp/go-plugin"
	"github.com/xchapter7x/lo"
)

var (
	cloudconfigs map[string]Record
	products     map[string]Record
)

func init() {
	cloudconfigs = make(map[string]Record)
	products = make(map[string]Record)
}

type Record struct {
	Name       string
	Path       string
	Properties map[string]interface{}
}

func ListCloudConfigs() map[string]Record {
	return cloudconfigs
}

func ListProducts() map[string]Record {
	return products
}

func RegisterProduct(pluginpath string) ([]pcli.Flag, error) {
	client, productPlugin := GetProductReference(pluginpath)
	defer client.Kill()
	meta := productPlugin.GetMeta()
	products[meta.Name] = Record{
		Name:       meta.Name,
		Path:       pluginpath,
		Properties: meta.Properties,
	}
	return productPlugin.GetFlags(), nil
}

func GetProductReference(pluginpath string) (*plugin.Client, product.ProductDeployer) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: product.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			product.PluginsMapHash: new(product.ProductPlugin),
		},
		Cmd: exec.Command(pluginpath, "plugin"),
	})
	rpcClient, err := client.Client()

	if err != nil {
		lo.G.Debug("we got an error:", err)
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense(product.PluginsMapHash)

	if err != nil {
		log.Fatal(err)
	}
	return client, raw.(product.ProductDeployer)
}

func RegisterCloudConfig(pluginpath string) ([]pcli.Flag, error) {
	client, ccPlugin := GetCloudConfigReference(pluginpath)
	defer client.Kill()
	meta := ccPlugin.GetMeta()
	cloudconfigs[meta.Name] = Record{
		Name:       meta.Name,
		Path:       pluginpath,
		Properties: meta.Properties,
	}
	return ccPlugin.GetFlags(), nil
}

func GetCloudConfigReference(pluginpath string) (*plugin.Client, cloudconfig.CloudConfigDeployer) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: cloudconfig.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			cloudconfig.PluginsMapHash: new(cloudconfig.CloudConfigPlugin),
		},
		Cmd: exec.Command(pluginpath, "plugin"),
	})

	rpcClient, err := client.Client()

	if err != nil {
		log.Fatal(err)
	}
	raw, err := rpcClient.Dispense(cloudconfig.PluginsMapHash)

	if err != nil {
		log.Fatal(err)
	}
	return client, raw.(cloudconfig.CloudConfigDeployer)
}
