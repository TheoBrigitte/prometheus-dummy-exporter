package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"

	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/config"
)

const (
	namespace = "dummy"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry").Default(":9510").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	configFile    = kingpin.Flag("config", "Path to config file").Default("").String()
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

func newCollector(namespace string, metrics []config.Metric) (*collector, error) {
	c := map[string]config.Metric{}
	counters := map[string]*prometheus.CounterVec{}
	gauges := map[string]*prometheus.GaugeVec{}
	for _, metric := range metrics {
		var keys []string
		for k := range metric.Labels {
			keys = append(keys, k)
		}
		keys = append([]string{"id"}, keys...)
		c[metric.Name] = metric
		switch metric.Type {
		case "counter":
			counters[metric.Name] = prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: namespace,
				Name:      metric.Name,
				Help:      "dummy counter",
			}, keys)
		case "gauge":
			gauges[metric.Name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      metric.Name,
				Help:      "dummy gauge",
			}, keys)
		default:
			return nil, fmt.Errorf("invalid type: %s for %s", metric.Type, metric.Name)
		}
	}
	return &collector{
		namespace: namespace,
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

func main() {
	kingpin.Version(version.Print("prometheus-dummy-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	conf := &config.Config{}
	if *configFile != "" {
		buf, err := ioutil.ReadFile(*configFile)
		if err != nil {
			log.Fatal("failed to read config file")
		}
		conf, err = config.Parse(buf)
		if err != nil {
			log.Fatal("invalid config format")
		}
	}

	collector, err := newCollector(namespace, conf.Metrics)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Infoln("listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
