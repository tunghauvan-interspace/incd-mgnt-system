package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/handlers"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/middleware"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize storage - use PostgreSQL if DATABASE_URL is provided, otherwise use memory store
	var store storage.Store
	var err error
	
	if cfg.DatabaseURL != "" {
		log.Println("Initializing PostgreSQL storage...")
		pgStore, err := storage.NewPostgresStore(cfg)
		if err != nil {
			log.Fatal("Failed to initialize PostgreSQL storage:", err)
		}
		store = pgStore
		log.Println("PostgreSQL storage initialized successfully")
	} else {
		log.Println("No database URL provided, using in-memory storage...")
		store, err = storage.NewMemoryStore()
		if err != nil {
			log.Fatal("Failed to initialize memory storage:", err)
		}
		log.Println("Memory storage initialized successfully")
	}

	// Initialize services
	metricsService := services.NewMetricsService()
	logger := services.NewLogger(cfg.LogLevel, true) // Use structured logging
	incidentService := services.NewIncidentService(store, metricsService)
	alertService := services.NewAlertService(store, incidentService, metricsService)
	notificationService := services.NewNotificationService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(incidentService, alertService, notificationService, metricsService, logger, store)

	// Setup middleware
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	
	// Apply middleware
	var h http.Handler = mux
	h = middleware.MetricsMiddleware(metricsService)(h)
	h = middleware.LoggingMiddleware(logger)(h)
	h = middleware.RequestIDMiddleware()(h)

	// Start background metrics updater
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			if err := incidentService.UpdatePrometheusMetrics(); err != nil {
				logger.Error("Failed to update Prometheus metrics", map[string]interface{}{
					"error": err.Error(),
				})
			}
			
			// Update database connection metrics if using PostgreSQL
			if pgStore, ok := store.(*storage.PostgresStore); ok {
				stats := pgStore.GetDBStats()
				metricsService.UpdateDBConnections(stats)
			}
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting Incident Management System", map[string]interface{}{
		"port": port,
		"log_level": cfg.LogLevel,
		"structured_logging": true,
		"metrics_enabled": cfg.MetricsEnabled,
	})

	if err := http.ListenAndServe(":"+port, h); err != nil {
		logger.Error("Server failed to start", map[string]interface{}{
			"error": err.Error(),
		})
		log.Fatal("Server failed to start:", err)
	}
}