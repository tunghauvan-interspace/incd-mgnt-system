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
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/middleware"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/models"
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
	userService          *services.UserService
	authService          *services.AuthService
	authHandler          *AuthHandler
}

// NewHandler creates a new handler
func NewHandler(
	incidentService *services.IncidentService,
	alertService *services.AlertService,
	notificationService *services.NotificationService,
	metricsService *services.MetricsService,
	logger *services.Logger,
	store storage.Store,
	userService *services.UserService,
	authService *services.AuthService,
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
	
	// Create auth handler
	authHandler := NewAuthHandler(userService, authService, logger)
	
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
		userService:         userService,
		authService:         authService,
		authHandler:         authHandler,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Authentication routes (public)
	mux.HandleFunc("/api/auth/register", h.authHandler.Register)
	mux.HandleFunc("/api/auth/login", h.authHandler.Login)
	mux.HandleFunc("/api/auth/refresh", h.authHandler.RefreshToken)
	
	// Protected authentication routes
	mux.HandleFunc("/api/auth/logout", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.authHandler.Logout)).ServeHTTP)
	mux.HandleFunc("/api/auth/profile", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.authHandler.GetProfile)).ServeHTTP)
	mux.HandleFunc("/api/auth/profile/update", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.authHandler.UpdateProfile)).ServeHTTP)
	mux.HandleFunc("/api/auth/password/change", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.authHandler.ChangePassword)).ServeHTTP)

	// API routes with rate limiting
	webhookHandler := ratelimit.WebhookRateLimitWrapper(h.rateLimitConfig, h.handleAlertmanagerWebhook)
	mux.HandleFunc("/api/webhooks/alertmanager", webhookHandler)
	
	// Protected API routes - require authentication
	mux.HandleFunc("/api/incidents/", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleIncidents)).ServeHTTP)
	mux.HandleFunc("/api/incidents", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleListIncidents)).ServeHTTP)
	mux.HandleFunc("/api/alerts", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleListAlerts)).ServeHTTP)
	mux.HandleFunc("/api/metrics", middleware.OptionalAuthMiddleware(h.authService)(http.HandlerFunc(h.handleGetMetrics)).ServeHTTP) // JSON metrics (deprecated)

	// Enhanced Incident Features - Protected API routes
	mux.HandleFunc("/api/incidents/search", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleIncidentSearch)).ServeHTTP)
	mux.HandleFunc("/api/incidents/bulk", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleIncidentBulkOperations)).ServeHTTP)
	mux.HandleFunc("/api/incidents/from-template", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleIncidentFromTemplate)).ServeHTTP)
	
	// Incident sub-resources (comments, tags, timeline, assignments)
	mux.HandleFunc("/api/incidents/", func(w http.ResponseWriter, r *http.Request) {
		middleware.AuthMiddleware(h.authService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathParts := strings.Split(r.URL.Path, "/")
			if len(pathParts) >= 5 {
				switch pathParts[4] {
				case "comments":
					h.handleIncidentComments(w, r)
				case "timeline":
					h.handleIncidentTimeline(w, r)
				case "tags":
					h.handleIncidentTags(w, r)
				case "assign":
					h.handleIncidentAssignment(w, r)
				default:
					h.handleIncidents(w, r) // Fallback to existing handler
				}
			} else {
				h.handleIncidents(w, r)
			}
		})).ServeHTTP(w, r)
	})

	// Template management
	mux.HandleFunc("/api/templates", middleware.AuthMiddleware(h.authService)(http.HandlerFunc(h.handleIncidentTemplates)).ServeHTTP)

	// Prometheus metrics endpoint (public for monitoring)
	mux.Handle("/metrics", promhttp.Handler())

	// Static files (CSS, JS, images, fonts, and other assets)
	mux.Handle("/css/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
	mux.Handle("/js/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
	mux.Handle("/images/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
	mux.Handle("/fonts/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
	mux.Handle("/assets/", http.StripPrefix("/", http.FileServer(http.Dir("web/static/"))))
	
	// SPA routes - serve Vue.js application for all non-API routes
	mux.HandleFunc("/", h.handleSPA)

	// Health check endpoints (public)
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/ready", h.handleReady)
	mux.HandleFunc("/db/stats", middleware.RequireRole(h.authService, "admin")(http.HandlerFunc(h.handleDBStats)).ServeHTTP)
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
		if r.Method == http.MethodPut {
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

// handleSPA serves the Vue.js Single Page Application
func (h *Handler) handleSPA(w http.ResponseWriter, r *http.Request) {
	// Check if it's an API route
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}
	
	// Only serve GET requests for the SPA
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Serve index.html for all frontend routes
	indexPath := filepath.Join("web/static", "index.html")
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, indexPath)
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

// handleDBStats handles database connection statistics requests
func (h *Handler) handleDBStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.logger.DebugWithRequest(r.Context(), "Database stats requested")

	stats := make(map[string]interface{})
	
	// Get database connection pool statistics
	if pgStore, ok := h.store.(*storage.PostgresStore); ok {
		dbStats := pgStore.GetDBStats()
		stats["database"] = map[string]interface{}{
			"type":                     "postgresql",
			"max_open_connections":     dbStats.MaxOpenConnections,
			"open_connections":         dbStats.OpenConnections,
			"in_use":                  dbStats.InUse,
			"idle":                    dbStats.Idle,
			"wait_count":              dbStats.WaitCount,
			"wait_duration":           dbStats.WaitDuration.String(),
			"max_idle_closed":         dbStats.MaxIdleClosed,
			"max_idle_time_closed":    dbStats.MaxIdleTimeClosed,
			"max_lifetime_closed":     dbStats.MaxLifetimeClosed,
		}
		
		// Add health status
		if err := pgStore.HealthCheck(); err != nil {
			stats["database"].(map[string]interface{})["health"] = "unhealthy"
			stats["database"].(map[string]interface{})["health_error"] = err.Error()
		} else {
			stats["database"].(map[string]interface{})["health"] = "healthy"
		}
	} else {
		stats["database"] = map[string]interface{}{
			"type":   "memory",
			"health": "healthy",
			"note":   "Memory store does not have connection statistics",
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"stats":     stats,
	})
}

// sendNotificationWithCircuitBreaker sends notifications with circuit breaker protection
func (h *Handler) sendNotificationWithCircuitBreaker(notificationFunc func() error) error {
	return h.circuitBreaker.Call(notificationFunc)
}

// Enhanced Incident Features - Comment Handlers

func (h *Handler) handleIncidentComments(w http.ResponseWriter, r *http.Request) {
	// Extract incident ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		h.writeErrorResponse(w, "Incident ID is required", http.StatusBadRequest)
		return
	}
	incidentID := pathParts[3]

	switch r.Method {
	case http.MethodGet:
		h.handleGetIncidentComments(w, r, incidentID)
	case http.MethodPost:
		h.handleAddIncidentComment(w, r, incidentID)
	default:
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetIncidentComments(w http.ResponseWriter, r *http.Request, incidentID string) {
	comments, err := h.incidentService.GetComments(incidentID)
	if err != nil {
		log.Printf("Failed to get comments for incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"comments": comments,
	})
}

func (h *Handler) handleAddIncidentComment(w http.ResponseWriter, r *http.Request, incidentID string) {
	var req struct {
		Content     string                     `json:"content"`
		CommentType models.IncidentCommentType `json:"comment_type"`
		UserID      string                     `json:"user_id"` // In real implementation, extract from auth
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		h.writeErrorResponse(w, "Content is required", http.StatusBadRequest)
		return
	}

	if req.CommentType == "" {
		req.CommentType = models.CommentTypeComment
	}

	if req.UserID == "" {
		req.UserID = "system" // Default for now
	}

	comment, err := h.incidentService.AddComment(incidentID, req.UserID, req.Content, req.CommentType, nil)
	if err != nil {
		log.Printf("Failed to add comment to incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *Handler) handleIncidentTimeline(w http.ResponseWriter, r *http.Request) {
	// Extract incident ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		h.writeErrorResponse(w, "Incident ID is required", http.StatusBadRequest)
		return
	}
	incidentID := pathParts[3]

	if r.Method != http.MethodGet {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	timeline, err := h.incidentService.GetTimeline(incidentID)
	if err != nil {
		log.Printf("Failed to get timeline for incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to retrieve timeline", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"timeline": timeline,
	})
}

// Enhanced Incident Features - Tag Handlers

func (h *Handler) handleIncidentTags(w http.ResponseWriter, r *http.Request) {
	// Extract incident ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		h.writeErrorResponse(w, "Incident ID is required", http.StatusBadRequest)
		return
	}
	incidentID := pathParts[3]

	switch r.Method {
	case http.MethodGet:
		h.handleGetIncidentTags(w, r, incidentID)
	case http.MethodPost:
		h.handleAddIncidentTags(w, r, incidentID)
	case http.MethodDelete:
		h.handleRemoveIncidentTags(w, r, incidentID)
	default:
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetIncidentTags(w http.ResponseWriter, r *http.Request, incidentID string) {
	tags, err := h.incidentService.GetTags(incidentID)
	if err != nil {
		log.Printf("Failed to get tags for incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to retrieve tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tags": tags,
	})
}

func (h *Handler) handleAddIncidentTags(w http.ResponseWriter, r *http.Request, incidentID string) {
	var req struct {
		Tags   []models.TemplateTag `json:"tags"`
		UserID string               `json:"user_id"` // In real implementation, extract from auth
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Tags) == 0 {
		h.writeErrorResponse(w, "At least one tag is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "system" // Default for now
	}

	err := h.incidentService.AddTags(incidentID, req.UserID, req.Tags)
	if err != nil {
		log.Printf("Failed to add tags to incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to add tags", http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, "Tags added successfully")
}

func (h *Handler) handleRemoveIncidentTags(w http.ResponseWriter, r *http.Request, incidentID string) {
	var req struct {
		TagNames []string `json:"tag_names"`
		UserID   string   `json:"user_id"` // In real implementation, extract from auth
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.TagNames) == 0 {
		h.writeErrorResponse(w, "At least one tag name is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "system" // Default for now
	}

	err := h.incidentService.RemoveTags(incidentID, req.UserID, req.TagNames)
	if err != nil {
		log.Printf("Failed to remove tags from incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to remove tags", http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, "Tags removed successfully")
}

// Enhanced Incident Features - Template Handlers

func (h *Handler) handleIncidentTemplates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListIncidentTemplates(w, r)
	case http.MethodPost:
		h.handleCreateIncidentTemplate(w, r)
	default:
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleListIncidentTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.incidentService.ListTemplates()
	if err != nil {
		log.Printf("Failed to list incident templates: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
	})
}

func (h *Handler) handleCreateIncidentTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.IncidentTemplate

	if err := json.NewDecoder(r.Body).Decode(&template); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if template.Name == "" {
		h.writeErrorResponse(w, "Template name is required", http.StatusBadRequest)
		return
	}

	if template.TitleTemplate == "" {
		h.writeErrorResponse(w, "Title template is required", http.StatusBadRequest)
		return
	}

	// Set defaults
	if template.CreatedBy == nil {
		userID := "system" // In real implementation, extract from auth
		template.CreatedBy = &userID
	}

	err := h.incidentService.CreateTemplate(&template)
	if err != nil {
		log.Printf("Failed to create incident template: %v", err)
		h.writeErrorResponse(w, "Failed to create template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

func (h *Handler) handleIncidentFromTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateIncidentFromTemplateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.TemplateID == "" {
		h.writeErrorResponse(w, "Template ID is required", http.StatusBadRequest)
		return
	}

	userID := "system" // In real implementation, extract from auth

	incident, err := h.incidentService.UseTemplate(&req, userID)
	if err != nil {
		log.Printf("Failed to create incident from template: %v", err)
		h.writeErrorResponse(w, "Failed to create incident from template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(incident)
}

// Enhanced Incident Features - Search Handler

func (h *Handler) handleIncidentSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.IncidentSearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100 // Maximum limit
	}

	response, err := h.incidentService.SearchIncidents(&req)
	if err != nil {
		log.Printf("Failed to search incidents: %v", err)
		h.writeErrorResponse(w, "Failed to search incidents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Enhanced Incident Features - Bulk Operations Handler

func (h *Handler) handleIncidentBulkOperations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.BulkOperationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.IncidentIDs) == 0 {
		h.writeErrorResponse(w, "At least one incident ID is required", http.StatusBadRequest)
		return
	}

	userID := "system" // In real implementation, extract from auth

	var response *models.BulkOperationResponse
	var err error

	switch req.Operation {
	case models.BulkOperationAcknowledge:
		assigneeID := "system" // Default assignee
		if assignee, ok := req.Parameters["assignee_id"].(string); ok {
			assigneeID = assignee
		}
		response, err = h.incidentService.BulkAcknowledge(req.IncidentIDs, assigneeID, userID)

	case models.BulkOperationUpdateStatus:
		statusStr, ok := req.Parameters["status"].(string)
		if !ok {
			h.writeErrorResponse(w, "Status parameter is required for status update operation", http.StatusBadRequest)
			return
		}
		status := models.IncidentStatus(statusStr)
		response, err = h.incidentService.BulkUpdateStatus(req.IncidentIDs, status, userID)

	default:
		h.writeErrorResponse(w, "Unsupported bulk operation", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Failed to perform bulk operation: %v", err)
		h.writeErrorResponse(w, "Failed to perform bulk operation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Enhanced Incident Features - Assignment Handler

func (h *Handler) handleIncidentAssignment(w http.ResponseWriter, r *http.Request) {
	// Extract incident ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 || pathParts[3] == "" {
		h.writeErrorResponse(w, "Incident ID is required", http.StatusBadRequest)
		return
	}
	incidentID := pathParts[3]

	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AssigneeID string `json:"assignee_id"`
		UserID     string `json:"user_id"` // In real implementation, extract from auth
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.AssigneeID == "" {
		h.writeErrorResponse(w, "Assignee ID is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "system" // Default for now
	}

	err := h.incidentService.AssignIncident(incidentID, req.AssigneeID, req.UserID)
	if err != nil {
		log.Printf("Failed to assign incident %s: %v", incidentID, err)
		h.writeErrorResponse(w, "Failed to assign incident", http.StatusInternalServerError)
		return
	}

	h.writeSuccessResponse(w, "Incident assigned successfully")
}

