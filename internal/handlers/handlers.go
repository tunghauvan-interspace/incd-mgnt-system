package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/circuitbreaker"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/idempotency"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/ratelimit"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/retry"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/validation"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler handles HTTP requests
type Handler struct {
	incidentService       *services.IncidentService
	alertService          *services.AlertService
	notificationService   *services.NotificationService
	webhookValidator      *validation.WebhookValidator
	idempotencyManager    *idempotency.WebhookIdempotencyManager
	retryer              *retry.Retryer
	rateLimitConfig      *ratelimit.RateLimitConfig
	circuitBreaker       *circuitbreaker.CircuitBreaker
	metricsService       *services.MetricsService
	logger               *services.Logger
	store                storage.Store
}

// NewHandler creates a new handler
func NewHandler(
	incidentService *services.IncidentService,
	alertService *services.AlertService,
	notificationService *services.NotificationService,
	metricsService *services.MetricsService,
	logger *services.Logger,
	store storage.Store,
) *Handler {
	// Initialize reliability components
	webhookValidator := validation.NewWebhookValidator()
	idempotencyStore := idempotency.NewMemoryIdempotencyStore()
	idempotencyManager := idempotency.NewWebhookIdempotencyManager(idempotencyStore, 10*time.Minute)
	
	// Create retry policy for webhook processing
	retryPolicy := &retry.RetryPolicy{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    2 * time.Second,
		Multiplier:  2.0,
	}
	retryer := retry.NewRetryer(retryPolicy, retry.DefaultIsRetryable)
	
	// Rate limiting configuration
	rateLimitConfig := ratelimit.DefaultWebhookRateLimit()
	
	// Circuit breaker for notification service
	cbConfig := circuitbreaker.DefaultConfig()
	cbConfig.OnStateChange = func(name string, from circuitbreaker.State, to circuitbreaker.State) {
		log.Printf("Circuit breaker '%s' changed state from %s to %s", name, from, to)
	}
	circuitBreaker := circuitbreaker.NewCircuitBreaker("notification-service", cbConfig)
	
	return &Handler{
		incidentService:     incidentService,
		alertService:        alertService,
		notificationService: notificationService,
		webhookValidator:    webhookValidator,
		idempotencyManager:  idempotencyManager,
		retryer:            retryer,
		rateLimitConfig:    rateLimitConfig,
		circuitBreaker:     circuitBreaker,
		metricsService:      metricsService,
		logger:              logger,
		store:               store,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// API routes with rate limiting
	webhookHandler := ratelimit.WebhookRateLimitWrapper(h.rateLimitConfig, h.handleAlertmanagerWebhook)
	mux.HandleFunc("/api/webhooks/alertmanager", webhookHandler)
	
	mux.HandleFunc("/api/incidents/", h.handleIncidents)
	mux.HandleFunc("/api/incidents", h.handleListIncidents)
	mux.HandleFunc("/api/alerts", h.handleListAlerts)
	mux.HandleFunc("/api/metrics", h.handleGetMetrics) // JSON metrics (deprecated)

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	
	// Dashboard
	mux.HandleFunc("/", h.handleDashboard)
	mux.HandleFunc("/incidents", h.handleIncidentsPage)
	mux.HandleFunc("/alerts", h.handleAlertsPage)

	// Health check endpoints
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/ready", h.handleReady)
}

// handleAlertmanagerWebhook handles incoming webhooks from Alertmanager with reliability improvements
func (h *Handler) handleAlertmanagerWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	// Validate HTTP method
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.InfoWithRequest(ctx, "Received Alertmanager webhook")

	// Validate HTTP method
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		h.writeErrorResponse(w, "Failed to read request body", http.StatusBadRequest)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Validate payload size (prevent large payload attacks)
	if len(body) > 1024*1024 { // 1MB limit
		log.Printf("Request body too large: %d bytes", len(body))
		h.writeErrorResponse(w, "Request body too large", http.StatusRequestEntityTooLarge)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Check for empty body
	if len(body) == 0 {
		h.writeErrorResponse(w, "Empty request body", http.StatusBadRequest)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Validate JSON schema
	if err := h.webhookValidator.ValidateAlertmanagerWebhook(body); err != nil {
		log.Printf("Webhook validation failed: %v", err)
		h.writeErrorResponse(w, fmt.Sprintf("Invalid webhook payload: %v", err), http.StatusBadRequest)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Check idempotency
	if h.idempotencyManager.IsAlreadyProcessed(body) {
		log.Printf("Duplicate webhook detected, returning cached response")
		h.writeSuccessResponse(w, "Duplicate request processed successfully")
		h.metricsService.RecordWebhookRequest("alertmanager", "success")
		return
	}

	// Parse webhook payload
	var webhook services.AlertmanagerWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		log.Printf("Failed to unmarshal webhook: %v", err)
		h.writeErrorResponse(w, "Invalid JSON structure", http.StatusBadRequest)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Process webhook with retry logic and circuit breaker
	err = h.retryer.Execute(ctx, func() error {
		return h.processWebhookWithCircuitBreaker(&webhook)
	})

	if err != nil {
		log.Printf("Failed to process webhook after retries: %v", err)
		h.writeErrorResponse(w, "Failed to process webhook", http.StatusInternalServerError)
		h.metricsService.RecordWebhookRequest("alertmanager", "error")
		return
	}

	// Mark as processed for idempotency
	if err := h.idempotencyManager.MarkAsProcessed(body); err != nil {
		log.Printf("Failed to mark webhook as processed: %v", err)
		// Don't fail the request for this error
	}

	// Return success response
	h.writeSuccessResponse(w, "Webhook processed successfully")
	h.metricsService.RecordWebhookRequest("alertmanager", "success")
	h.logger.InfoWithRequest(ctx, "Successfully processed Alertmanager webhook", map[string]interface{}{
		"alerts_count": len(webhook.Alerts),
		"status":       webhook.Status,
	})
}

// processWebhookWithCircuitBreaker processes webhook with circuit breaker protection
func (h *Handler) processWebhookWithCircuitBreaker(webhook *services.AlertmanagerWebhook) error {
	return h.circuitBreaker.Call(func() error {
		return h.alertService.ProcessAlertmanagerWebhook(webhook)
	})
}

// writeErrorResponse writes a structured error response
func (h *Handler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := map[string]interface{}{
		"status": "error",
		"error":  message,
		"code":   statusCode,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write error response: %v", err)
	}
}

// writeSuccessResponse writes a structured success response
func (h *Handler) writeSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to write success response: %v", err)
	}
}

// handleIncidents handles incident-related requests with different methods and paths
func (h *Handler) handleIncidents(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	
	if path == "/api/incidents" || path == "/api/incidents/" {
		if r.Method == http.MethodGet {
			h.handleListIncidents(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}
	
	// Parse incident ID from path /api/incidents/{id} or /api/incidents/{id}/action
	pathParts := strings.Split(strings.TrimPrefix(path, "/api/incidents/"), "/")
	if len(pathParts) == 0 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	
	incidentID := pathParts[0]
	
	if len(pathParts) == 1 {
		// /api/incidents/{id}
		if r.Method == http.MethodGet {
			h.handleGetIncident(w, r, incidentID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if len(pathParts) == 2 {
		// /api/incidents/{id}/action
		action := pathParts[1]
		if r.Method == http.MethodPost {
			switch action {
			case "acknowledge":
				h.handleAcknowledgeIncident(w, r, incidentID)
			case "resolve":
				h.handleResolveIncident(w, r, incidentID)
			default:
				http.Error(w, "Unknown action", http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Invalid path", http.StatusBadRequest)
	}
}

// handleListIncidents returns all incidents
func (h *Handler) handleListIncidents(w http.ResponseWriter, r *http.Request) {
	incidents, err := h.incidentService.ListIncidents()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incidents)
}

// handleGetIncident returns a specific incident
func (h *Handler) handleGetIncident(w http.ResponseWriter, r *http.Request, id string) {
	incident, err := h.incidentService.GetIncident(id)
	if err != nil {
		http.Error(w, "Incident not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident)
}

// AcknowledgeIncidentRequest represents the request to acknowledge an incident
type AcknowledgeIncidentRequest struct {
	AssigneeID string `json:"assignee_id"`
}

// handleAcknowledgeIncident acknowledges an incident
func (h *Handler) handleAcknowledgeIncident(w http.ResponseWriter, r *http.Request, id string) {
	var req AcknowledgeIncidentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.incidentService.AcknowledgeIncident(id, req.AssigneeID); err != nil {
		http.Error(w, "Failed to acknowledge incident", http.StatusInternalServerError)
		return
	}

	// Get updated incident
	incident, err := h.incidentService.GetIncident(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send notification with circuit breaker
	if err := h.sendNotificationWithCircuitBreaker(func() error {
		return h.notificationService.NotifyIncidentAcknowledged(incident)
	}); err != nil {
		log.Printf("Failed to send acknowledgment notification: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident)
}

// handleResolveIncident resolves an incident
func (h *Handler) handleResolveIncident(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.incidentService.ResolveIncident(id); err != nil {
		http.Error(w, "Failed to resolve incident", http.StatusInternalServerError)
		return
	}

	// Get updated incident
	incident, err := h.incidentService.GetIncident(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send notification with circuit breaker
	if err := h.sendNotificationWithCircuitBreaker(func() error {
		return h.notificationService.NotifyIncidentResolved(incident)
	}); err != nil {
		log.Printf("Failed to send resolution notification: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incident)
}

// handleListAlerts returns all alerts
func (h *Handler) handleListAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	alerts, err := h.alertService.ListAlerts()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// handleGetMetrics returns incident metrics
func (h *Handler) handleGetMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	metrics, err := h.incidentService.CalculateMetrics()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handleDashboard serves the main dashboard page
func (h *Handler) handleDashboard(w http.ResponseWriter, r *http.Request) {
	log.Printf("Dashboard handler called with method: %s, path: %s", r.Method, r.URL.Path)
	
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.serveTemplate(w, r, "dashboard.html")
}

// handleIncidentsPage serves the incidents page
func (h *Handler) handleIncidentsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.serveTemplate(w, r, "incidents.html")
}

// handleAlertsPage serves the alerts page
func (h *Handler) handleAlertsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.serveTemplate(w, r, "alerts.html")
}

// handleHealth handles health check requests
func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.DebugWithRequest(r.Context(), "Health check requested")
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// handleReady handles readiness probe requests
func (h *Handler) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.DebugWithRequest(r.Context(), "Readiness check requested")

	// Check database connectivity
	ready := true
	checks := make(map[string]interface{})

	// Test database connection
	if pgStore, ok := h.store.(*storage.PostgresStore); ok {
		if err := pgStore.HealthCheck(); err != nil {
			ready = false
			checks["database"] = map[string]interface{}{
				"status": "unhealthy",
				"error":  err.Error(),
			}
			h.logger.ErrorWithRequest(r.Context(), "Database health check failed", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			checks["database"] = map[string]interface{}{
				"status": "healthy",
			}
		}
	} else {
		// Memory store is always ready
		checks["database"] = map[string]interface{}{
			"status": "healthy",
			"type":   "memory",
		}
	}

	status := "ready"
	statusCode := http.StatusOK
	if !ready {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"checks":    checks,
	})
}

// sendNotificationWithCircuitBreaker sends notifications with circuit breaker protection
func (h *Handler) sendNotificationWithCircuitBreaker(notificationFunc func() error) error {
	return h.circuitBreaker.Call(notificationFunc)
}

// serveTemplate serves HTML templates
func (h *Handler) serveTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	templatePath := filepath.Join("web/templates", templateName)
	
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, templatePath)
}