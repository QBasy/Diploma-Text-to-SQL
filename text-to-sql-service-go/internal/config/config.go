package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TTSQL struct {
		APIKey string `yaml:"API_KEY"`
		PORT   string `yaml:"PORT"`
		GROC   struct {
			APIKey  string `yaml:"API_KEY"`
			MODEL   string `yaml:"MODEL"`
			BASEURL string `yaml:"BASE_URL"`
		} `yaml:"GROC"`
	} `yaml:"TTSQL"`
}

func LoadConfig(filename string) (*Config, error) {
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
