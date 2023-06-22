package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerHost            string `yaml:"ServerHost"`
	ServerPort            int    `yaml:"ServerPort"`
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

	return config, nil
}
