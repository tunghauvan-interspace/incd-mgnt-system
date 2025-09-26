package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
)

// Handler handles HTTP requests
type Handler struct {
	incidentService     *services.IncidentService
	alertService        *services.AlertService
	notificationService *services.NotificationService
}

// NewHandler creates a new handler
func NewHandler(
	incidentService *services.IncidentService,
	alertService *services.AlertService,
	notificationService *services.NotificationService,
) *Handler {
	return &Handler{
		incidentService:     incidentService,
		alertService:        alertService,
		notificationService: notificationService,
	}
}

// RegisterRoutes registers all HTTP routes
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("/api/webhooks/alertmanager", h.handleAlertmanagerWebhook)
	mux.HandleFunc("/api/incidents/", h.handleIncidents)
	mux.HandleFunc("/api/incidents", h.handleListIncidents)
	mux.HandleFunc("/api/alerts", h.handleListAlerts)
	mux.HandleFunc("/api/metrics", h.handleGetMetrics)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))
	
	// Dashboard
	mux.HandleFunc("/", h.handleDashboard)
	mux.HandleFunc("/incidents", h.handleIncidentsPage)
	mux.HandleFunc("/alerts", h.handleAlertsPage)

	// Health check
	mux.HandleFunc("/health", h.handleHealth)
}

// handleAlertmanagerWebhook handles incoming webhooks from Alertmanager
func (h *Handler) handleAlertmanagerWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webhook services.AlertmanagerWebhook
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.alertService.ProcessAlertmanagerWebhook(&webhook); err != nil {
		log.Printf("Error processing alertmanager webhook: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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

	// Send notification
	if err := h.notificationService.NotifyIncidentAcknowledged(incident); err != nil {
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

	// Send notification
	if err := h.notificationService.NotifyIncidentResolved(incident); err != nil {
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// serveTemplate serves HTML templates
func (h *Handler) serveTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	templatePath := filepath.Join("web/templates", templateName)
	
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, templatePath)
}