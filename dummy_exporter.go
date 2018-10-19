package main

import (
	"fmt"
	"github.com/kobtea/dummy_exporter/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"net/http"
)

const (
	namespace = "dummy"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry").Default(":9999").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	configFile    = kingpin.Flag("config", "Path to config file").Default("").String()
)

type collector struct {
	namespace string
	config    map[string]config.Metric
	desc      map[string][]*prometheus.Desc
}

func newCollector(namespace string, metrics []config.Metric) (*collector, error) {
	c := map[string]config.Metric{}
	d := map[string][]*prometheus.Desc{}
	for _, metric := range metrics {
		var keys []string
		for k := range metric.Labels {
			keys = append(keys, k)
		}
		var descs []*prometheus.Desc
		for i := 0; i < metric.Size; i++ {
			descs = append(descs, prometheus.NewDesc(fmt.Sprintf("%s_%s_%d", namespace, metric.Name, i), "dummy", keys, nil))
		}

		c[metric.Name] = metric
		d[metric.Name] = descs
	}
	return &collector{
		namespace: namespace,
		config:    c,
		desc:      d,
	}, nil
}

func (collector collector) Describe(ch chan<- *prometheus.Desc) {
	for _, descs := range collector.desc {
		for _, desc := range descs {
			ch <- desc
		}
	}
}

func (collector collector) Collect(ch chan<- prometheus.Metric) {
	for _, descs := range collector.desc {
		for _, desc := range descs {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1)
		}
	}
}

func main() {
	kingpin.Version(version.Print("dummy_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	buf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal("failed to read config file")
	}
	conf, err := config.Parse(buf)
	if err != nil {
		log.Fatal("invalid config format")
	}

	collector, err := newCollector(namespace, conf.Metrics)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Dummy Exporter</title></head>
             <body>
             <h1>Dummy Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Infoln("listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
