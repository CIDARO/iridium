package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
}

func GetConfig(configPath string) (*Config, error) {
	config := &Config{}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env != "test" && env != "production" && env != "development" {
		return nil, errors.New("environment must be either 'test', 'production' or 'development'")
	}

	file, err := os.Open(configPath + "/" + env)
	if err != nil {
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
