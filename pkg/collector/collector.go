package collector

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/config"
)

type collector struct {
	namespace string
	config    map[string]config.Metric
	counters  map[string]*prometheus.CounterVec
	gauges    map[string]*prometheus.GaugeVec
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(conf *config.Config) (*collector, error) {
	c := map[string]config.Metric{}
	counters := map[string]*prometheus.CounterVec{}
	gauges := map[string]*prometheus.GaugeVec{}
	for _, metric := range conf.Metrics {
		var keys []string
		for k := range metric.Labels {
			keys = append(keys, k)
		}
		keys = append([]string{"id"}, keys...)
		c[metric.Name] = metric
		switch metric.Type {
		case "counter":
			counters[metric.Name] = prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: conf.Namespace,
				Name:      metric.Name,
				Help:      "dummy counter",
			}, keys)
		case "gauge":
			gauges[metric.Name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: conf.Namespace,
				Name:      metric.Name,
				Help:      "dummy gauge",
			}, keys)
		default:
			return nil, fmt.Errorf("invalid type: %s for %s", metric.Type, metric.Name)
		}
	}
	return &collector{
		namespace: conf.Namespace,
		config:    c,
		counters:  counters,
		gauges:    gauges,
	}, nil
}

func (collector collector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range collector.counters {
		metric.Describe(ch)
	}
	for _, metric := range collector.gauges {
		metric.Describe(ch)
	}
}

func (collector collector) Collect(ch chan<- prometheus.Metric) {
	for name, conf := range collector.config {
		for i := 0; i < conf.Size; i++ {
			labels := map[string]string{"id": strconv.Itoa(i)}
			for key, vals := range conf.Labels {
				labels[key] = vals[i%len(vals)]
			}
			switch conf.Type {
			case "counter":
				collector.counters[name].With(labels).Inc()
				collector.counters[name].With(labels).Collect(ch)
			case "gauge":
				collector.gauges[name].With(labels).Set(rand.Float64())
				collector.gauges[name].With(labels).Collect(ch)
			default:
				log.Errorf("invalid type: %s for %s", conf.Type, conf.Name)
			}
		}
	}
}
