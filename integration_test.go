package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/handlers"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/middleware"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

// TestIntegration_MetricsAndMonitoring tests the complete monitoring infrastructure
func TestIntegration_MetricsAndMonitoring(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test environment
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	// Initialize services with monitoring
	metricsService := services.NewMetricsService()
	logger := services.NewLogger("debug", true)
	incidentService := services.NewIncidentService(store, metricsService)
	alertService := services.NewAlertService(store, incidentService, metricsService)
	notificationService := services.NewNotificationService(&config.Config{})

	// Initialize handlers
	handler := handlers.NewHandler(incidentService, alertService, notificationService, metricsService, logger, store)

	// Setup routes with middleware
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Apply middleware
	var h http.Handler = mux
	h = middleware.MetricsMiddleware(metricsService)(h)
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RequestIDMiddleware()(h)

	server := httptest.NewServer(h)
	defer server.Close()

	t.Run("PrometheusMetricsEndpoint", func(t *testing.T) {
		testPrometheusMetricsEndpoint(t, server)
	})

	t.Run("HealthAndReadinessEndpoints", func(t *testing.T) {
		testHealthAndReadinessEndpoints(t, server)
	})

	t.Run("HTTPRequestInstrumentation", func(t *testing.T) {
		testHTTPRequestInstrumentation(t, server)
	})

	t.Run("WebhookProcessingInstrumentation", func(t *testing.T) {
		testWebhookProcessingInstrumentation(t, server)
	})

	t.Run("StructuredLogging", func(t *testing.T) {
		testStructuredLogging(t, server)
	})

	t.Run("RequestIDTracing", func(t *testing.T) {
		testRequestIDTracing(t, server)
	})

	t.Run("MetricsCollection", func(t *testing.T) {
		testMetricsCollection(t, server, incidentService)
	})
}

func testPrometheusMetricsEndpoint(t *testing.T, server *httptest.Server) {
	// Make some requests first to populate metrics
	http.Get(server.URL + "/health")
	http.Get(server.URL + "/api/incidents")
	
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	bodyStr := string(body)

	// Check for essential Prometheus metrics (only the ones that should always be present)
	expectedMetrics := []string{
		"http_requests_total",
		"http_request_duration_seconds", 
		"http_requests_in_flight",
		// Skip metrics that might not be populated yet
		// "db_query_duration_seconds",
		// "incidents_total",
		// "alerts_total", 
		// "webhook_requests_total",
		"incident_mtta_seconds",
		"incident_mttr_seconds",
	}

	for _, metric := range expectedMetrics {
		if !strings.Contains(bodyStr, metric) {
			t.Errorf("Expected metric %s not found in response", metric)
		}
	}

	// Verify Prometheus format
	if !strings.Contains(bodyStr, "# HELP") || !strings.Contains(bodyStr, "# TYPE") {
		t.Error("Response does not appear to be in Prometheus format")
	}
}

func testHealthAndReadinessEndpoints(t *testing.T, server *httptest.Server) {
	// Test health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to get health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for health, got %d", resp.StatusCode)
	}

	var healthResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if status, ok := healthResp["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", healthResp["status"])
	}

	if _, ok := healthResp["timestamp"]; !ok {
		t.Error("Expected timestamp in health response")
	}

	// Test readiness endpoint
	resp, err = http.Get(server.URL + "/ready")
	if err != nil {
		t.Fatalf("Failed to get ready endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for ready, got %d", resp.StatusCode)
	}

	var readyResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&readyResp); err != nil {
		t.Fatalf("Failed to decode ready response: %v", err)
	}

	if status, ok := readyResp["status"].(string); !ok || status != "ready" {
		t.Errorf("Expected status 'ready', got %v", readyResp["status"])
	}

	if checks, ok := readyResp["checks"].(map[string]interface{}); !ok || len(checks) == 0 {
		t.Error("Expected checks in ready response")
	}
}

func testHTTPRequestInstrumentation(t *testing.T, server *httptest.Server) {
	// Make several requests to generate metrics
	endpoints := []string{
		"/api/incidents",
		"/api/alerts",
		"/health",
		"/ready",
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get(server.URL + endpoint)
		if err != nil {
			t.Errorf("Failed to get %s: %v", endpoint, err)
			continue
		}
		resp.Body.Close()
	}

	// Check that metrics are recorded
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read metrics: %v", err)
	}

	bodyStr := string(body)

	// Verify that HTTP requests are tracked
	for _, endpoint := range endpoints {
		if !strings.Contains(bodyStr, fmt.Sprintf(`path="%s"`, endpoint)) {
			t.Errorf("Expected to find metrics for path %s", endpoint)
		}
	}

	// Verify status codes are tracked
	if !strings.Contains(bodyStr, `status_code="200"`) {
		t.Error("Expected to find status_code label in metrics")
	}
}

func testWebhookProcessingInstrumentation(t *testing.T, server *httptest.Server) {
	// Create a test webhook payload
	webhook := map[string]interface{}{
		"status": "firing",
		"alerts": []map[string]interface{}{
			{
				"fingerprint": "test-fingerprint-123",
				"status":      "firing",
				"labels": map[string]interface{}{
					"alertname": "TestAlert",
					"severity":  "critical",
				},
				"annotations": map[string]interface{}{
					"summary":     "Test alert for integration testing",
					"description": "This is a test alert",
				},
			},
		},
	}

	payload, err := json.Marshal(webhook)
	if err != nil {
		t.Fatalf("Failed to marshal webhook: %v", err)
	}

	// Send webhook
	resp, err := http.Post(
		server.URL+"/api/webhooks/alertmanager",
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		t.Fatalf("Failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for webhook, got %d", resp.StatusCode)
	}

	// Check that webhook metrics are recorded
	metricsResp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer metricsResp.Body.Close()

	body, err := io.ReadAll(metricsResp.Body)
	if err != nil {
		t.Fatalf("Failed to read metrics: %v", err)
	}

	bodyStr := string(body)

	// Verify webhook metrics
	if !strings.Contains(bodyStr, `webhook_requests_total{source="alertmanager",status="success"}`) {
		t.Error("Expected webhook success metric")
	}

	// Verify incident creation metrics
	if !strings.Contains(bodyStr, `incidents_total{severity="critical",status="open"}`) {
		t.Error("Expected incident creation metric")
	}
}

func testStructuredLogging(t *testing.T, server *httptest.Server) {
	// Capture stdout to check for structured logs
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Make a request to trigger logging
	resp, err := http.Get(server.URL + "/api/incidents")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	resp.Body.Close()

	// Give some time for logs to be written
	time.Sleep(100 * time.Millisecond)

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout

	logOutput, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read log output: %v", err)
	}

	logStr := string(logOutput)

	// Verify structured logging format (JSON)
	if !strings.Contains(logStr, `"timestamp":`) {
		t.Error("Expected timestamp field in structured logs")
	}
	if !strings.Contains(logStr, `"level":`) {
		t.Error("Expected level field in structured logs")
	}
	if !strings.Contains(logStr, `"message":`) {
		t.Error("Expected message field in structured logs")
	}
	if !strings.Contains(logStr, `"request_id":`) {
		t.Error("Expected request_id field in structured logs")
	}
}

func testRequestIDTracing(t *testing.T, server *httptest.Server) {
	// Make a request and check for request ID header
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	requestID := resp.Header.Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set")
	}

	// Verify it's a valid UUID format
	if len(requestID) != 36 {
		t.Errorf("Expected UUID format (36 chars), got %d chars: %s", len(requestID), requestID)
	}
	if !strings.Contains(requestID, "-") {
		t.Errorf("Expected UUID format with hyphens, got: %s", requestID)
	}
}

func testMetricsCollection(t *testing.T, server *httptest.Server, incidentService *services.IncidentService) {
	// Update Prometheus metrics
	err := incidentService.UpdatePrometheusMetrics()
	if err != nil {
		t.Fatalf("Failed to update Prometheus metrics: %v", err)
	}

	// Get metrics
	resp, err := http.Get(server.URL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read metrics: %v", err)
	}

	bodyStr := string(body)

	// Verify MTTA and MTTR metrics are present (even if zero)
	if !strings.Contains(bodyStr, "incident_mtta_seconds") {
		t.Error("Expected MTTA metric to be present")
	}
	if !strings.Contains(bodyStr, "incident_mttr_seconds") {
		t.Error("Expected MTTR metric to be present")
	}

	// Verify incidents by status metrics
	if !strings.Contains(bodyStr, "incidents_by_status") {
		t.Error("Expected incidents by status metrics")
	}
}

// TestIntegration_DatabaseMonitoring tests database-specific monitoring features
func TestIntegration_DatabaseMonitoring(t *testing.T) {
	// Skip if no database URL provided
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database monitoring integration test")
	}

	cfg := &config.Config{
		DatabaseURL:       dbURL,
		DBMaxOpenConns:    10,
		DBMaxIdleConns:    2,
		DBConnMaxLifetime: 5 * time.Minute,
		LogLevel:          "debug",
		MetricsEnabled:    true,
	}

	// Initialize with PostgreSQL
	pgStore, err := storage.NewPostgresStore(cfg)
	if err != nil {
		t.Fatalf("Failed to create PostgreSQL store: %v", err)
	}
	defer pgStore.Close()
	
	// Use interface
	var store storage.Store = pgStore

	// Initialize services
	metricsService := services.NewMetricsService()
	logger := services.NewLogger(cfg.LogLevel, true)
	incidentService := services.NewIncidentService(store, metricsService)
	alertService := services.NewAlertService(store, incidentService, metricsService)
	notificationService := services.NewNotificationService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(incidentService, alertService, notificationService, metricsService, logger, store)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Apply middleware
	var h http.Handler = mux
	h = middleware.MetricsMiddleware(metricsService)(h)
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RequestIDMiddleware()(h)

	server := httptest.NewServer(h)
	defer server.Close()

	t.Run("DatabaseHealthCheck", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/ready")
		if err != nil {
			t.Fatalf("Failed to get ready endpoint: %v", err)
		}
		defer resp.Body.Close()

		var readyResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&readyResp); err != nil {
			t.Fatalf("Failed to decode ready response: %v", err)
		}

		checks, ok := readyResp["checks"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected checks in ready response")
		}

		dbCheck, ok := checks["database"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected database check")
		}

		if status, ok := dbCheck["status"].(string); !ok || status != "healthy" {
			t.Errorf("Expected database status 'healthy', got %v", dbCheck["status"])
		}
	})

	t.Run("DatabaseConnectionMetrics", func(t *testing.T) {
		// Update database metrics directly using pgStore
		stats := pgStore.GetDBStats()
		metricsService.UpdateDBConnections(stats)

		// Get metrics
		resp, err := http.Get(server.URL + "/metrics")
		if err != nil {
			t.Fatalf("Failed to get metrics: %v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		// Verify database connection metrics
		expectedDBMetrics := []string{
			`db_connections{status="open"}`,
			`db_connections{status="idle"}`,
			`db_connections{status="in_use"}`,
		}

		for _, metric := range expectedDBMetrics {
			if !strings.Contains(bodyStr, metric) {
				t.Errorf("Expected database metric %s", metric)
			}
		}
	})
}

// TestIntegration_EndToEndMonitoring tests the complete end-to-end monitoring flow
func TestIntegration_EndToEndMonitoring(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping end-to-end integration test in short mode")
	}

	// Setup in-memory environment
	store, err := storage.NewMemoryStore()
	if err != nil {
		t.Fatalf("Failed to create memory store: %v", err)
	}

	cfg := &config.Config{
		LogLevel:       "info",
		MetricsEnabled: true,
		Port:          "0", // Let test server choose port
	}

	// Initialize services
	metricsService := services.NewMetricsService()
	logger := services.NewLogger(cfg.LogLevel, true)
	incidentService := services.NewIncidentService(store, metricsService)
	alertService := services.NewAlertService(store, incidentService, metricsService)
	notificationService := services.NewNotificationService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(incidentService, alertService, notificationService, metricsService, logger, store)

	// Setup routes with full middleware stack
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	var h http.Handler = mux
	h = middleware.MetricsMiddleware(metricsService)(h)
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RequestIDMiddleware()(h)

	server := httptest.NewServer(h)
	defer server.Close()

	// Test complete workflow: webhook -> incident creation -> metrics
	t.Run("CompleteWorkflow", func(t *testing.T) {
		// 1. Send webhook to create incident
		webhook := map[string]interface{}{
			"status": "firing",
			"alerts": []map[string]interface{}{
				{
					"fingerprint": "e2e-test-fingerprint",
					"status":      "firing",
					"labels": map[string]interface{}{
						"alertname": "E2ETestAlert",
						"severity":  "high",
					},
					"annotations": map[string]interface{}{
						"summary": "End-to-end test alert",
					},
				},
			},
		}

		payload, _ := json.Marshal(webhook)
		resp, err := http.Post(
			server.URL+"/api/webhooks/alertmanager",
			"application/json",
			bytes.NewBuffer(payload),
		)
		if err != nil {
			t.Fatalf("Failed to send webhook: %v", err)
		}
		resp.Body.Close()

		// 2. Verify incident was created
		incidentsResp, err := http.Get(server.URL + "/api/incidents")
		if err != nil {
			t.Fatalf("Failed to get incidents: %v", err)
		}
		defer incidentsResp.Body.Close()

		var incidents []map[string]interface{}
		if err := json.NewDecoder(incidentsResp.Body).Decode(&incidents); err != nil {
			t.Fatalf("Failed to decode incidents: %v", err)
		}

		if len(incidents) == 0 {
			t.Fatal("Expected at least one incident to be created")
		}

		// 3. Update Prometheus metrics
		if err := incidentService.UpdatePrometheusMetrics(); err != nil {
			t.Fatalf("Failed to update metrics: %v", err)
		}

		// 4. Verify all metrics are properly updated
		metricsResp, err := http.Get(server.URL + "/metrics")
		if err != nil {
			t.Fatalf("Failed to get metrics: %v", err)
		}
		defer metricsResp.Body.Close()

		body, err := io.ReadAll(metricsResp.Body)
		if err != nil {
			t.Fatalf("Failed to read metrics: %v", err)
		}

		bodyStr := string(body)

		// Verify business metrics
		if !strings.Contains(bodyStr, `incidents_total{severity="high",status="open"} 1`) {
			t.Error("Expected incident creation metric")
		}

		if !strings.Contains(bodyStr, `webhook_requests_total{source="alertmanager",status="success"} 1`) {
			t.Error("Expected webhook success metric")
		}

		// Verify HTTP metrics for all our requests
		if !strings.Contains(bodyStr, `http_requests_total{method="POST",path="/api/webhooks/alertmanager",status_code="200"}`) {
			t.Error("Expected webhook HTTP metric")
		}

		if !strings.Contains(bodyStr, `http_requests_total{method="GET",path="/api/incidents",status_code="200"}`) {
			t.Error("Expected incidents HTTP metric")
		}

		if !strings.Contains(bodyStr, `http_requests_total{method="GET",path="/metrics",status_code="200"}`) {
			t.Error("Expected metrics HTTP metric")
		}
	})
}