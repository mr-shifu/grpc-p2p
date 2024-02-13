package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Peer struct {
	Addr        string `yaml:"addr"`
	Attributes  map[string]string
	Name        string `yaml:"name"`
	ClusterName string `yaml:"clusterName"`
}

type Config struct {
	Local     Peer   `yaml:"local"`
	Bootstrap []Peer `yaml:"bootstrap"`
}

func FromFile(path string) (*Config, error) {
	// read connection profile yaml file
	cpb, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal yaml to connection profile struct
	cp := &Config{}
	if err := yaml.Unmarshal(cpb, cp); err != nil {
		return nil, err
	}
	return cp, nil
}
