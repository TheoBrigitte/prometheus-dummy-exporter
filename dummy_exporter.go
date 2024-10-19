package main

import (
	"fmt"
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
	logLevel      = kingpin.Flag("log.level", fmt.Sprintf("Log level: %v", logLevelsString())).Default(log.InfoLevel.String()).String()
)

func logLevelsString() (levels []string) {
	for _, level := range log.AllLevels {
		levels = append(levels, level.String())
	}
	return
}

func main() {
	kingpin.Version(version.Print("prometheus-dummy-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %v", err)
	}
	log.SetLevel(level)

	conf := config.New()
	if *configFile != "" {
		err := conf.ReadFromFile(*configFile)
		if err != nil {
			log.Fatalf("failed to read config file: %v", err)
		}
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
	log.Fatal(http.ListenAndServe(*listenAddress, nil)) // nolint:gosec
}
