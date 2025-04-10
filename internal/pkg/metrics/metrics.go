package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	// Технические метрики
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec

	// Бизнес-метрики
	PVZCreated       prometheus.Counter
	ReceptionCreated prometheus.Counter
	ProductAdded     prometheus.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		RequestCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method", "endpoint"},
		),
		PVZCreated: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "pvz_created_total",
				Help: "Total number of created PVZs",
			},
		),
		ReceptionCreated: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "reception_created_total",
				Help: "Total number of created receptions",
			},
		),
		ProductAdded: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "product_added_total",
				Help: "Total number of added products",
			},
		),
	}
}
