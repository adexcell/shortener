package main

import (
	"context"
	"log"
	"net/http"

	"github.com/adexcell/shortener.git/internal/config"
	"github.com/adexcell/shortener.git/internal/migrations"
	"github.com/adexcell/shortener.git/internal/repository/postgres"
	"github.com/adexcell/shortener.git/internal/service"
	httpTransport "github.com/adexcell/shortener.git/internal/transport/http"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting service on %s", cfg.HTTPServer)

	log.Println("Applying migrations...")
	if err := migrations.RunMigrations(cfg.DatabaseDSN); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	log.Println("Migrations applied successfully")

	ctx := context.Background()
	store, err := postgres.NewStorage(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to storage: %v", err)
	}
	defer store.Close()

	svc := service.NewService(store)

	handler := httpTransport.NewHandler(svc)

	srv := &http.Server{
		Addr:    cfg.HTTPServer,
		Handler: handler.InitRoutes(),
	}

	log.Printf("Starting HTTP server on %s", cfg.HTTPServer)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
