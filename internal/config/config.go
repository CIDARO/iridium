package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/CIDARO/iridium/internal/utils"

	"gopkg.in/yaml.v2"
)

// Configuration config struct
type Configuration struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Memcache struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}
	Backends    []string `yaml:"backends"`
	MaxAttempts int      `yaml:"max_attempts"`
	MaxRetries  int      `yaml:"max_retries"`
	Metrics     bool     `yaml:"metrics"`
}

// Config is the global configuration object
var Config *Configuration = &Configuration{}

// GetConfig initialized the configuration
func GetConfig(configPath string, vars ...string) {
	// Retrieve ENV environment variable
	env := os.Getenv("ENV")
	// If vars is provided, the first one is the override for the env
	if len(vars) > 0 {
		env = vars[0]
	}
	// Check if no env is provided
	if env == "" {
		// If it is not provided, set it as development
		env = "development"
	}
	// Check if env is either 'test', 'production' or 'development'
	if env != "test" && env != "production" && env != "development" {
		// Raise an error if it does not match the criteria
		utils.Logger.Panic(errors.New("environment must be either 'test', 'production' or 'development'"))
	}
	// Open the configuration file
	file, err := os.Open(fmt.Sprintf("%s/%s.yml", configPath, env))
	if err != nil {
		// Return an error if it's given
		utils.Logger.Panic(err)
	}
	// Decode the yaml file into the struct
	yamlDecoder := yaml.NewDecoder(file)
	if err := yamlDecoder.Decode(&Config); err != nil {
		utils.Logger.Panic(err)
	}
}

// ValidatePath validates the input path
func ValidatePath(path string) (*string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("'%s' is not a valid directory", path)
	}
	return &path, nil
}
