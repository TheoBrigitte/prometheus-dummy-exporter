package config

import (
	"os"

	"sigs.k8s.io/yaml"
)

var (
	defaultConfig = Config{
		Namespace: "dummy",
	}
)

type Config struct {
	Metrics   []Metric `yaml:"metrics"`
	Namespace string   `yaml:"namespace"`
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

	var c = defaultConfig
	if err := yaml.Unmarshal(buf, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
