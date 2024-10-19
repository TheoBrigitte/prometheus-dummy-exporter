package main

import (
	"strconv"
	"testing"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/collector"
	"github.com/TheoBrigitte/prometheus-dummy-exporter/pkg/config"
)

func TestCollect(t *testing.T) {
	conf := config.New()

	c, err := collector.New(conf)
	if err != nil {
		t.Fatal(err)
	}

	ch := make(chan prometheus.Metric)
	go func(ch chan prometheus.Metric) {
		for range ch {
		}
	}(ch)

	c.Collect(ch)

	close(ch)
}

func BenchmarkCollect(b *testing.B) {
	alice := config.Metric{
		Name:   "alice",
		Type:   config.MetricTypeGauge,
		Size:   500000,
		Labels: map[string][]string{},
	}
	addLabels(&alice)

	bob := config.Metric{
		Name:   "bob",
		Type:   config.MetricTypeCounter,
		Size:   500000,
		Labels: map[string][]string{},
	}
	addLabels(&bob)

	conf := &config.Config{
		Metrics: []config.Metric{
			alice, bob,
		},
	}

	c, err := collector.New(conf)
	if err != nil {
		b.Fatal(err)
	}

	ch := make(chan prometheus.Metric)
	go func(ch chan prometheus.Metric) {
		for range ch {
		}
	}(ch)

	b.ResetTimer()
	for range b.N {
		c.Collect(ch)
	}

	close(ch)
}

func addLabels(m *config.Metric) {
	for i := 0; i < 10; i++ {
		is := strconv.Itoa(i)
		key := "l" + is
		for j := 0; j < 10; j++ {
			m.Labels[key] = append(m.Labels[key], "l"+is+strconv.Itoa(j))
		}
	}
}
