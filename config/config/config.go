package config

import (
	"os"

	"github.com/vandi37/vanerrors"
	"gopkg.in/yaml.v3"
)

// The errors
const (
	ErrorReadingConfig      = "error reading config"
	ErrorUnmarshalingConfig = "error unmarshaling config"
)

// The database connection config
type DatabaseCfg struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// The application config
type AppConfig struct {
	IsService bool   `yaml:"is_service"`
	Duration  string `yaml:"duration"`
}

// The standard config
type Config struct {
	Port     int         `yaml:"port"`
	Database DatabaseCfg `yaml:"database"`
	App      AppConfig   `yaml:"app"`
	Salt     string      `yaml:"salt"`
	Key      string      `yaml:"key"`
}

// Loads config from the yaml file
func LoadConfig(path string) (*Config, error) {
	// Getting the config data
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorReadingConfig, err, vanerrors.EmptyHandler)

	}

	// unmarshal the config
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorUnmarshalingConfig, err, vanerrors.EmptyHandler)
	}

	return &config, nil
}
