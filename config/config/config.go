package config

import (
	"log"
	"os"

	"github.com/VandiKond/vanerrors"
	"gopkg.in/yaml.v3"
)

type Config interface {
	GetDBConData() DBCfg
}

type DBCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type StdCfg struct {
	Port     string `yaml:"port"`
	Database DBCfg  `yaml:"database"`
}

func (cfg StdCfg) GetDBConData() DBCfg {
	return cfg.Database
}

func LoadConfig(path string) *StdCfg {
	var config StdCfg

	data, err := os.ReadFile(path)
	if err != nil {
		err = vanerrors.NewWrap("error reading config file", err, vanerrors.EmptyHandler)
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		err = vanerrors.NewWrap("error unmarshaling config file", err, vanerrors.EmptyHandler)
		log.Fatal(err)
	}

	return &config
}
