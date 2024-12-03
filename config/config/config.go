package config

import (
	"os"

	"github.com/VandiKond/vanerrors"
	"gopkg.in/yaml.v3"
)

type Config interface {
	GetDBConData() DBCfg
}

type DBCfg struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type StdCfg struct {
	Port     string `yaml:"port"`
	Database DBCfg  `yaml:"database"`
	Salt     string `yaml:"salt"`
}

func (cfg StdCfg) GetDBConData() DBCfg {
	return cfg.Database
}

func LoadConfig(path string) (*StdCfg, error) {
	var config StdCfg

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, vanerrors.NewWrap("error reading config file", err, vanerrors.EmptyHandler)

	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, vanerrors.NewWrap("error unmarshaling config file", err, vanerrors.EmptyHandler)
	}

	return &config, nil
}
