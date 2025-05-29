package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Environment struct {
	AccountId         string `json:"account_id,omitempty"`
	Cloud             string `json:"cloud"`
	Name              string `json:"name"`
	NetworkId         string `json:"network_id,omitempty"`
	ProjectId         string `json:"project,omitempty"`
	Region            string `json:"region,omitempty"`
	StorageExternalId string `json:"storage_external_id,omitempty"`
	SubscriptionId    string `json:"subscription_id,omitempty"`
	Url               string `json:"url,omitempty"`
	VpcId             string `json:"vpc_id,omitempty"`
}

var environments []Environment

func LoadEnvironments(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "failed to read environments file")
	}

	if err := json.Unmarshal(data, &environments); err != nil {
		return errors.Wrap(err, "failed to parse environments")
	}

	return nil
}

func GetEnvironment(name string) (*Environment, error) {
	for _, env := range environments {
		if env.Name == name {
			return &env, nil
		}
	}
	return nil, errors.New("environment not found")
}

func init() {
	if err := LoadEnvironments("environments.json"); err != nil {
		panic(err)
	}
}
