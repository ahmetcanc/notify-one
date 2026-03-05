// ...existing code...
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ahmetcanc/notify-one/api"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/cache"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/config"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/database"
	"github.com/ahmetcanc/notify-one/internal/usecase"
	"github.com/ahmetcanc/notify-one/migrations"
)

func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode)

	// Run database migrations
	if err := migrations.RunMigrations(dsn); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize PostgreSQL pool
	dbPool, err := database.NewConnection(ctx, cfg.ToDatabaseConfig())
	if err != nil {
		log.Fatalf("PostgreSQL connection failed: %v", err)
	}
	defer dbPool.Close()

	// Initialize Redis connection
	rdb, err := cache.NewRedisConnection(ctx, cfg.ToRedisConfig())
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer rdb.Close()

	// Initialize Repositories
	notificationRepo := database.NewPostgresNotificationRepository(dbPool)
	// Logger for check
	log.Printf("Repository initialized: %T", notificationRepo)

	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, rdb)

	notificationHandler := api.NewNotificationHandler(notificationUsecase)
	// Routes
	http.HandleFunc("POST /api/v1/notifications", notificationHandler.Send)
	http.HandleFunc("POST /api/v1/notifications/batch", notificationHandler.SendBatch)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverAddr := ":" + cfg.AppPort
	log.Printf("🚀 Server starting on %s", serverAddr)

	// Start HTTP server
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
