package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all application metrics
type Metrics struct {
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	GRPCRequestsTotal     *prometheus.CounterVec
	GRPCRequestDuration   *prometheus.HistogramVec
	DatabaseQueriesTotal  *prometheus.CounterVec
	DatabaseQueryDuration *prometheus.HistogramVec
	EventsPublishedTotal  *prometheus.CounterVec
	EventsConsumedTotal   *prometheus.CounterVec
	ActiveConnections     prometheus.Gauge
}

// NewMetrics creates and registers all metrics
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request latency in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		GRPCRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "grpc_requests_total",
				Help:      "Total number of gRPC requests",
			},
			[]string{"method", "status"},
		),
		GRPCRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "grpc_request_duration_seconds",
				Help:      "gRPC request latency in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method"},
		),
		DatabaseQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "database_queries_total",
				Help:      "Total number of database queries",
			},
			[]string{"operation", "status"},
		),
		DatabaseQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "database_query_duration_seconds",
				Help:      "Database query latency in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"operation"},
		),
		EventsPublishedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "events_published_total",
				Help:      "Total number of events published",
			},
			[]string{"event_type", "status"},
		),
		EventsConsumedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "events_consumed_total",
				Help:      "Total number of events consumed",
			},
			[]string{"event_type", "status"},
		),
		ActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "active_connections",
				Help:      "Number of active connections",
			},
		),
	}
}
