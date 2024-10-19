package collector

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/config"
)

type collector struct {
	counters []counter
	gauges   []gauge
}

type gauge struct {
	config config.Metric
	vec    *prometheus.GaugeVec
}

type counter struct {
	config config.Metric
	vec    *prometheus.CounterVec
}

func New(conf *config.Config) (*collector, error) {
	counters := []counter{}
	gauges := []gauge{}

	for _, metric := range conf.Metrics {
		keys := slices.Collect(maps.Keys(metric.Labels))
		keys = append([]string{"id"}, keys...)

		switch metric.Type {
		case config.MetricTypeCounter:
			c := counter{
				config: metric,
				vec: prometheus.NewCounterVec(prometheus.CounterOpts{
					Namespace: conf.Namespace,
					Name:      metric.Name,
					Help:      "dummy counter",
				}, keys),
			}
			counters = append(counters, c)
		case config.MetricTypeGauge:
			g := gauge{
				config: metric,
				vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Namespace: conf.Namespace,
					Name:      metric.Name,
					Help:      "dummy gauge",
				}, keys),
			}
			gauges = append(gauges, g)
		default:
			return nil, fmt.Errorf("invalid type: %s for %s", metric.Type, metric.Name)
		}
	}

	c := &collector{
		counters: counters,
		gauges:   gauges,
	}
	return c, nil
}

func (collector collector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range collector.counters {
		c.vec.Describe(ch)
	}
	for _, g := range collector.gauges {
		g.vec.Describe(ch)
	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, counter := range c.counters {
		for i := 0; i < counter.config.Size; i++ {
			labels := counter.config.GenerateLabels(i)
			counter.vec.With(labels).Inc()
			counter.vec.With(labels).Collect(ch)
		}
	}

	for _, gauge := range c.gauges {
		for i := 0; i < gauge.config.Size; i++ {
			labels := gauge.config.GenerateLabels(i)
			gauge.vec.With(labels).Set(rand.Float64()) //nolint:gosec
			gauge.vec.With(labels).Collect(ch)
		}
	}
}
