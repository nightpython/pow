package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerHost            string `yaml:"ServerHost" envconfig:"SERVER_HOST"`
	ServerPort            int    `yaml:"ServerPort" envconfig:"SERVER_PORT"`
	HashcashZerosCount    int    `yaml:"HashcashZerosCount"`
	HashcashMaxIterations int    `yaml:"HashcashMaxIterations"`
}

// Load loads the configuration from the specified YAML file path.
func Load(filePath string) (*Config, error) {
	// Read the YAML file
	yamlData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the YAML data into the Config struct
	config := &Config{}
	err = yaml.Unmarshal(yamlData, config)
	if err != nil {
		return nil, err
	}
	err = envconfig.Process("", config)

	return config, err
}
