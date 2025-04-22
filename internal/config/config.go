package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ApiGatewayURL        string `yaml:"ApiGatewayURL"`
	ApiKey               string `yaml:"ApiKey"`
	PORT                 string `yaml:"PORT"`
	VisualisationService string `yaml:"VisualisationService"`
	TextToSQLService     string `yaml:"TextToSQLService"`
}

func GetConfig(filename string) (*Config, error) {
	var config *Config
	file, err := os.Open(filename)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("error decoding config file: %v", err)
	}
	return config, nil
}
