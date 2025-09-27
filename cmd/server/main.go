package main

import (
	"log"
	"net/http"
	"os"

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

	log.Printf("Starting Incident Management System on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}