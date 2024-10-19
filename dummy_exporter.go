package main

import (
	"net/http"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"

	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/collector"
	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/config"
)

var (
	listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry").Default(":9510").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	configFile    = kingpin.Flag("config", "Path to config file").Default("").String()
)

func main() {
	kingpin.Version(version.Print("prometheus-dummy-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	conf, err := config.NewFromFile(*configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	collector, err := collector.New(conf)
	if err != nil {
		log.Fatalf("failed to create collector: %v", err)
	}

	err = prometheus.Register(collector)
	if err != nil {
		log.Fatalf("failed to register collector: %v", err)
	}

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	log.Infoln("listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
