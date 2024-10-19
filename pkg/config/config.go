package config

import (
	"encoding/json"
	"os"
	"strconv"

	"sigs.k8s.io/yaml"
)

var (
	defaultConfig = Config{
		Namespace: "dummy",
		Metrics: []Metric{
			{
				Name: "http_requests",
				Type: MetricTypeGauge,
				Size: 10,
				Labels: map[string][]string{
					"code":   {"100", "200", "300", "400", "500"},
					"method": {"GET", "POST", "PUT", "DELETE"},
				},
			},
		},
	}
)

type Config struct {
	Metrics   []Metric `yaml:"metrics"`
	Namespace string   `yaml:"namespace,omitempty"`
}

type Metric struct {
	Name   string       `yaml:"name"`
	Type   MetricType   `yaml:"type"`
	Size   int          `yaml:"size"`
	Labels MetricLabels `yaml:"labels,omitempty"`
}

type MetricLabels map[string][]string

func New() *Config {
	return &defaultConfig
}

func (c *Config) ReadFromFile(configFile string) error {
	buf, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return c.Parse(buf)
}

func (c *Config) Parse(buf []byte) error {
	if err := yaml.Unmarshal(buf, c); err != nil {
		return err
	}

	return nil
}

func (m *Metric) GenerateLabels(i int) map[string]string {
	labels := map[string]string{"id": strconv.Itoa(i)}
	for key, vals := range m.Labels {
		labels[key] = vals[i%len(vals)]
	}

	return labels
}

func (m *MetricLabels) UnmarshalJSON(data []byte) (err error) {
	labels := make(map[string][]string)

	err = json.Unmarshal(data, &labels)
	if err != nil {
		return err
	}

	*m = labels

	return nil
}
