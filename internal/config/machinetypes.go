package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type MachineTypeOsConfig struct {
	OsType         string `json:"os_type"`
	OsDistribution string `json:"os_distribution"`
	OsDiskSizeGB   int    `json:"os_disk_size_gb"`
}

type MachineTypeCloudConfig struct {
	InstanceType string `json:"instance_type"`
	Image        string `json:"image"`
}

type MachineTypeConfig struct {
	Os     MachineTypeOsConfig               `json:"os"`
	Clouds map[string]MachineTypeCloudConfig `json:"clouds"`
}

var MachineTypesConfig map[string]MachineTypeConfig

func LoadMachineTypes(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "failed to read machinetypes file")
	}
	return json.Unmarshal(data, &MachineTypesConfig)
}

func GetMachineTypeConfig(machineType, cloud string) (*MachineTypeCloudConfig, *MachineTypeOsConfig, error) {
	mt, ok := MachineTypesConfig[machineType]
	if !ok {
		return nil, nil, errors.Errorf("machinetype %s not found", machineType)
	}
	conf, ok := mt.Clouds[cloud]
	if !ok {
		return nil, nil, errors.Errorf("cloud %s not found for machinetype %s", cloud, machineType)
	}
	return &conf, &mt.Os, nil
}

func init() {
	_ = LoadMachineTypes("machinetypes.json")
}
