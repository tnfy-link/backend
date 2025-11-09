package stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	redirectsTotal prometheus.Counter
}

func newMetrics() *metrics {
	return &metrics{
		redirectsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Subsystem: "stats",
				Name:      "redirects_total",
				Help:      "Total number of link redirects",
			},
		),
	}
}

func (m *metrics) IncRedirects() {
	m.redirectsTotal.Inc()
}
