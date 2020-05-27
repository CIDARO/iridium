package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Target url.URL
	DbPath string `yaml:"db_path"`
}

func GetConfig(configPath string, target string) (*Config, error) {
	config := &Config{}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	if env != "test" && env != "production" && env != "development" {
		return nil, errors.New("environment must be either 'test', 'production' or 'development'")
	}

	file, err := os.Open(fmt.Sprintf("%s/%s.yml", configPath, env))
	if err != nil {
		return nil, err
	}

	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(target)
	if err != nil {
		return config, nil
	}

	config.Target = *parsedURL

	return config, nil
}
