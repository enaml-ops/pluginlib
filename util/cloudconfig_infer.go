package pluginutil

import (
	"fmt"
	"strings"

	"github.com/enaml-ops/enaml"
)

type CloudConfigInferer interface {
	InferDefaultVMType() string
	InferDefaultDiskType() string
	InferDefaultNetwork() string
	InferDefaultAZ() string
}

type CloudConfigInfer struct {
	CloudConfig enaml.CloudConfigManifest
}

//NewCloudConfigInferFromBytes - get a cloud configinferer from a cloud config
//bytes array
func NewCloudConfigInferFromBytes(b []byte) *CloudConfigInfer {
	c := enaml.NewCloudConfigManifest(b)
	return &CloudConfigInfer{
		CloudConfig: *c,
	}
}

//NewCloudConfigInfer - get a cloud configinferer from a cloud config enaml
//object
func NewCloudConfigInfer(c enaml.CloudConfigManifest) *CloudConfigInfer {
	return &CloudConfigInfer{
		CloudConfig: c,
	}
}

func (s *CloudConfigInfer) InferDefaultVMType() string {
	var name = ""

	if len(s.CloudConfig.VMTypes) > 0 {
		name = s.CloudConfig.VMTypes[0].Name
	}
	return name
}

func (s *CloudConfigInfer) InferDefaultDiskType() string {
	var name = ""

	if len(s.CloudConfig.DiskTypes) > 0 {
		name = s.CloudConfig.DiskTypes[0].Name
	}
	return name
}

func (s *CloudConfigInfer) InferDefaultNetwork() string {
	var name = ""

	if len(s.CloudConfig.Networks) > 0 {
		name = s.CloudConfig.Networks[0].(map[interface{}]interface{})["name"].(string)
	}
	return name
}

func (s *CloudConfigInfer) InferDefaultAZ() string {
	var name = ""
	var names []string

	if len(s.CloudConfig.Networks) > 0 {
		for _, az := range s.CloudConfig.AZs {

			fmt.Println(az)
			names = append(names, az.Name)
		}
		name = strings.Join(names, ",")
	}
	return name
}
