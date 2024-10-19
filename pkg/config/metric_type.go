package config

import (
	"encoding/json"
	"fmt"
)

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
)

type MetricType uint8

func (mt MetricType) String() string {
	switch mt {
	case MetricTypeCounter:
		return "counter"
	case MetricTypeGauge:
		return "gauge"
	}

	return fmt.Sprintf("MetricType(%d)", mt)
}

func (mt *MetricType) UnmarshalJSON(data []byte) (err error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if *mt, err = parseMetricType(value); err != nil {
		return err
	}

	return nil
}

func parseMetricType(value string) (MetricType, error) {
	switch value {
	case "counter":
		return MetricTypeCounter, nil
	case "gauge":
		return MetricTypeGauge, nil
	}

	return MetricType(0), fmt.Errorf("%q is not a valid MetricTypemetric", value)
}
