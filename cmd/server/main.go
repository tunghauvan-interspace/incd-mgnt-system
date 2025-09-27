package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/config"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/handlers"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/middleware"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/services"
	"github.com/tunghauvan-interspace/incd-mgnt-system/internal/storage"
)

func main() {
	// Load and validate configuration
	cfg, err := config.LoadAndValidateConfig()
	if err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	log.Printf("Starting Incident Management System...")
	log.Printf("Configuration loaded successfully")
	
	if cfg.DebugMode {
		log.Printf("DEBUG MODE: Configuration details:")
		log.Printf("  - Port: %s", cfg.Port)
		log.Printf("  - Log Level: %s", cfg.LogLevel)
		log.Printf("  - Database: %s", func() string {
			if cfg.DatabaseURL != "" {
				return "PostgreSQL (configured)"
			}
			return "In-Memory"
		}())
		log.Printf("  - Notifications: %t", cfg.HasNotificationConfigured())
		log.Printf("  - TLS: %t", cfg.IsTLSEnabled())
		log.Printf("  - Metrics: %t (port %s)", cfg.MetricsEnabled, cfg.MetricsPort)
	}

	// Initialize storage - use PostgreSQL if DATABASE_URL is provided, otherwise use memory store
	var store storage.Store
	
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
		log.Println("WARNING: Using in-memory storage - data will be lost on restart")
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

	// Create HTTP server with configured timeouts
	server := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        mux,
		ReadTimeout:    cfg.ServerReadTimeout,
		WriteTimeout:   cfg.ServerWriteTimeout,
		IdleTimeout:    cfg.ServerIdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Info("Starting Incident Management System", map[string]interface{}{
		"port": cfg.Port,
		"log_level": cfg.LogLevel,
		"structured_logging": true,
		"metrics_enabled": cfg.MetricsEnabled,
	})

	// Start server in a goroutine
	go func() {
		if cfg.IsTLSEnabled() {
			log.Printf("Starting HTTPS server on port %s", cfg.Port)
			if err := server.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server failed to start: %v", err)
			}
		} else {
			log.Printf("Starting HTTP server on port %s", cfg.Port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server failed to start: %v", err)
			}
		}
	}()

	log.Printf("Incident Management System started successfully")
	if !cfg.HasNotificationConfigured() {
		log.Printf("WARNING: No notification methods configured - alerts will not be sent")
	}

	// Wait for interrupt signal
	<-ctx.Done()
	stop()
	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown gracefully")
	}

	// Close storage connections
	if pgStore, ok := store.(*storage.PostgresStore); ok {
		pgStore.Close()
		log.Println("Database connections closed")
	}

	log.Println("Shutdown complete")
}