package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsInterface interface {
	IncPVZCreated()
	IncReceptionCreated()
	IncProductAdded()
	ObserveRequestDuration(method, endpoint string, duration float64)
	IncRequestCount(method, endpoint, status string)
}

type Metrics struct {
	RequestCount    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec

	PVZCreated       prometheus.Counter
	ReceptionCreated prometheus.Counter
	ProductAdded     prometheus.Counter
}

func NewMetrics() *Metrics {
	return &Metrics{
		RequestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5},
			},
			[]string{"method", "endpoint"},
		),
		PVZCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "pvz_created_total",
				Help: "Total number of created PVZs",
			},
		),
		ReceptionCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "reception_created_total",
				Help: "Total number of created receptions",
			},
		),
		ProductAdded: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "product_added_total",
				Help: "Total number of added products",
			},
		),
	}
}

func (m *Metrics) IncPVZCreated() {
	m.PVZCreated.Inc()
}

func (m *Metrics) IncReceptionCreated() {
	m.ReceptionCreated.Inc()
}

func (m *Metrics) IncProductAdded() {
	m.ProductAdded.Inc()
}

func (m *Metrics) ObserveRequestDuration(method, endpoint string, duration float64) {
	m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func (m *Metrics) IncRequestCount(method, endpoint, status string) {
	m.RequestCount.WithLabelValues(method, endpoint, status).Inc()
}
