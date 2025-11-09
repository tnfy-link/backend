package links

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsErrorReason string

const (
	metricsErrorReasonValidationFailed   metricsErrorReason = "validation_failed"
	metricsErrorReasonIDGenerationFailed metricsErrorReason = "id_generation_failed"
	metricsErrorReasonStorageFailed      metricsErrorReason = "storage_failed"
)

type metrics struct {
	createdTotal prometheus.Counter
	errorsTotal  *prometheus.CounterVec
}

func newMetrics() *metrics {
	return &metrics{
		createdTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Subsystem: "links",
				Name:      "created_total",
				Help:      "Total number of shortened links created",
			},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: "links",
				Name:      "errors_total",
				Help:      "Total number of link errors",
			},
			[]string{"reason"},
		),
	}
}

func (m *metrics) IncCreated() {
	m.createdTotal.Inc()
}

func (m *metrics) IncError(reason metricsErrorReason) {
	m.errorsTotal.WithLabelValues(string(reason)).Inc()
}
