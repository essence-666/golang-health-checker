package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CheckerData struct {
	Retries uint     `yaml:"retries"`
	Timeout uint     `yaml:"timeout"`
	Period  uint     `yaml:"period"`
	URLs    []string `yaml:"urls"`
}

type CheckerConfig struct {
	ConfigData CheckerData `yaml:"checker"`
}

func ReadConfig() (*CheckerData, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	var config CheckerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config.ConfigData, nil
}
