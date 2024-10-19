package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Metrics []Metric `yaml:"metrics"`
}

type Metric struct {
	Name   string              `yaml:"name"`
	Type   string              `yaml:"type"`
	Size   int                 `yaml:"size"`
	Labels map[string][]string `yaml:"labels"`
}

func NewFromFile(configFile string) (*Config, error) {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	return Parse(buf)
}

func Parse(buf []byte) (*Config, error) {
	var c Config
	if err := yaml.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
