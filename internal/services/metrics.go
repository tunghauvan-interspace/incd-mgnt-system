package services

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsService handles Prometheus metrics collection
type MetricsService struct {
	// HTTP metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge

	// Database metrics
	dbQueryDuration *prometheus.HistogramVec
	dbConnections   *prometheus.GaugeVec

	// Business metrics
	incidentsTotal    *prometheus.CounterVec
	alertsTotal       *prometheus.CounterVec
	incidentsByStatus *prometheus.GaugeVec
	mtta              prometheus.Gauge
	mttr              prometheus.Gauge

	// Webhook metrics
	webhookRequestsTotal *prometheus.CounterVec
	notificationsSent    *prometheus.CounterVec
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	return &MetricsService{
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status_code"},
		),
		httpRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Current number of HTTP requests being served",
			},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "db_query_duration_seconds",
				Help: "Database query duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"query_type", "table"},
		),
		dbConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "db_connections",
				Help: "Current database connections",
			},
			[]string{"status"}, // open, idle, in_use
		),
		incidentsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "incidents_total",
				Help: "Total number of incidents created",
			},
			[]string{"severity", "status"},
		),
		alertsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "alerts_total",
				Help: "Total number of alerts processed",
			},
			[]string{"status"},
		),
		incidentsByStatus: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "incidents_by_status",
				Help: "Current number of incidents by status",
			},
			[]string{"status", "severity"},
		),
		mtta: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "incident_mtta_seconds",
				Help: "Mean Time To Acknowledge in seconds",
			},
		),
		mttr: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "incident_mttr_seconds",
				Help: "Mean Time To Resolve in seconds",
			},
		),
		webhookRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "webhook_requests_total",
				Help: "Total number of webhook requests processed",
			},
			[]string{"source", "status"},
		),
		notificationsSent: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "notifications_sent_total",
				Help: "Total number of notifications sent",
			},
			[]string{"channel", "status"},
		),
	}
}

// RecordHTTPRequest records an HTTP request metric
func (m *MetricsService) RecordHTTPRequest(method, path, statusCode string, duration time.Duration) {
	m.httpRequestsTotal.WithLabelValues(method, path, statusCode).Inc()
	m.httpRequestDuration.WithLabelValues(method, path, statusCode).Observe(duration.Seconds())
}

// IncrementHTTPInFlight increments the in-flight HTTP requests counter
func (m *MetricsService) IncrementHTTPInFlight() {
	m.httpRequestsInFlight.Inc()
}

// DecrementHTTPInFlight decrements the in-flight HTTP requests counter
func (m *MetricsService) DecrementHTTPInFlight() {
	m.httpRequestsInFlight.Dec()
}

// RecordDBQuery records a database query metric
func (m *MetricsService) RecordDBQuery(queryType, table string, duration time.Duration) {
	m.dbQueryDuration.WithLabelValues(queryType, table).Observe(duration.Seconds())
}

// UpdateDBConnections updates database connection metrics
func (m *MetricsService) UpdateDBConnections(stats sql.DBStats) {
	m.dbConnections.WithLabelValues("open").Set(float64(stats.OpenConnections))
	m.dbConnections.WithLabelValues("idle").Set(float64(stats.Idle))
	m.dbConnections.WithLabelValues("in_use").Set(float64(stats.InUse))
}

// RecordIncidentCreated records a new incident creation
func (m *MetricsService) RecordIncidentCreated(severity, status string) {
	m.incidentsTotal.WithLabelValues(severity, status).Inc()
}

// UpdateIncidentsByStatus updates the current incidents by status gauge
func (m *MetricsService) UpdateIncidentsByStatus(status, severity string, count float64) {
	m.incidentsByStatus.WithLabelValues(status, severity).Set(count)
}

// UpdateMTTA updates the Mean Time To Acknowledge metric
func (m *MetricsService) UpdateMTTA(mtta time.Duration) {
	m.mtta.Set(mtta.Seconds())
}

// UpdateMTTR updates the Mean Time To Resolve metric
func (m *MetricsService) UpdateMTTR(mttr time.Duration) {
	m.mttr.Set(mttr.Seconds())
}

// RecordAlertProcessed records an alert processing event
func (m *MetricsService) RecordAlertProcessed(status string) {
	m.alertsTotal.WithLabelValues(status).Inc()
}

// RecordWebhookRequest records a webhook request
func (m *MetricsService) RecordWebhookRequest(source, status string) {
	m.webhookRequestsTotal.WithLabelValues(source, status).Inc()
}

// RecordNotificationSent records a notification sending event
func (m *MetricsService) RecordNotificationSent(channel, status string) {
	m.notificationsSent.WithLabelValues(channel, status).Inc()
}