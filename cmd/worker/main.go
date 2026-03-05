package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/cache"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/config"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/database"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	// Connect to PostgreSQL
	dbPool, err := database.NewConnection(ctx, cfg.ToDatabaseConfig())
	if err != nil {
		log.Fatalf("PostgreSQL connection failed: %v", err)
	}
	defer dbPool.Close()

	// Connect to Redis
	rdb, err := cache.NewRedisConnection(ctx, cfg.ToRedisConfig())
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer rdb.Close()

	repo := database.NewPostgresNotificationRepository(dbPool)

	log.Println("👷 Worker started, listening for notifications...")

	for {
		// Wait for notification ID from Redis
		result, err := rdb.PopFromQueue(ctx, "notification_queue")
		if err != nil {
			log.Printf("Queue error: %v", err)
			continue
		}

		// result[0] is the key (queue name), result[1] is the value (notification ID)
		notificationID := result[1]
		log.Printf("Processing notification: %s", notificationID)

		// 1. Fetch details from Database
		n, err := repo.GetByID(ctx, notificationID)
		if err != nil {
			log.Printf("Failed to fetch notification %s: %v", notificationID, err)
			continue
		}

		// 2. Simulate sending process
		fmt.Printf(">>> SENDING [%s] to %s: %s\n", n.Channel, n.Recipient, n.Content)
		time.Sleep(2 * time.Second) // Simulate network delay

		// 3. Update Status in DB
		err = repo.UpdateStatus(ctx, notificationID, domain.StatusSent)
		if err != nil {
			log.Printf("Failed to update status for %s: %v", notificationID, err)
			continue
		}

		log.Printf("✅ Notification %s marked as SENT", notificationID)
	}
}
