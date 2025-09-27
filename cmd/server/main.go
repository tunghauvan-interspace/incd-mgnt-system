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
	incidentService := services.NewIncidentService(store)
	alertService := services.NewAlertService(store, incidentService)
	notificationService := services.NewNotificationService(cfg)

	// Initialize handlers
	handler := handlers.NewHandler(incidentService, alertService, notificationService)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server with timeouts for better reliability
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel to listen for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Start server in goroutine
	go func() {
		log.Printf("Starting Incident Management System on port %s", port)
		log.Printf("Server configured with:")
		log.Printf("  - Read timeout: %v", server.ReadTimeout)
		log.Printf("  - Write timeout: %v", server.WriteTimeout)
		log.Printf("  - Idle timeout: %v", server.IdleTimeout)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			stop <- syscall.SIGTERM
		}
	}()

	// Wait for interrupt signal
	<-stop

	log.Println("Shutting down server gracefully...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}

	// Cleanup storage connections if needed
	if pgStore, ok := store.(*storage.PostgresStore); ok {
		log.Println("Closing database connections...")
		if err := pgStore.Close(); err != nil {
			log.Printf("Error closing database connections: %v", err)
		} else {
			log.Println("Database connections closed successfully")
		}
	}

	log.Println("Application shutdown complete")
}