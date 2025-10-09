package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics.
	HttpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed, labeled by status, method and path.",
	}, []string{"method", "path", "status"})

	HttpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of request durations.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})

	HttpActiveRequests = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_active_requests",
		Help: "Current number of active HTTP requests.",
	}, []string{"method", "path"})

	// Rate Limit metrics.
	RateLimitBlockedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "rate_limit_blocked_total",
		Help: "Number of requests blocked by the rate limiter.",
	}, []string{"method", "path"})

	// Config reload metrics.
	ConfigReloads = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "config_reload_total",
		Help: "Number of configuration reloads attempted, labeled by status.",
	}, []string{"status"})

	ConfigReloadLast = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "config_reload_last_timestamp_seconds",
		Help: "Timestamp of last configuration reload attempt.",
	})

	// Business / app metrics examples.
	JobsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "jobs_processed_total",
		Help: "Total number of background jobs processed successfully.",
	}, []string{"type"})

	ErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "errors_total",
		Help: "Total number of errors encountered, labeled by source.",
	}, []string{"source"})
)

func RecordConfigReload(success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	ConfigReloads.WithLabelValues(status).Inc()
	ConfigReloadLast.SetToCurrentTime()
}
