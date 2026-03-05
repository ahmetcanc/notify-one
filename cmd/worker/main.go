package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ahmetcanc/notify-one/internal/domain"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/cache"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/config"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/database"
	"github.com/ahmetcanc/notify-one/internal/infrastructure/provider"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	// Connect to PostgreSQL
	dbPool, err := database.NewConnection(ctx, cfg.ToDatabaseConfig())
	if err != nil {
		log.Fatalf("Postgres connection failed: %v", err)
	}
	defer dbPool.Close()

	// Connect to Redis
	rdb, err := cache.NewRedisConnection(ctx, cfg.ToRedisConfig())
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
	}
	defer rdb.Close()

	repo := database.NewPostgresNotificationRepository(dbPool)

	// Initialize Webhook Provider as requested in assessment
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		webhookURL = "https://webhook.site/YOUR-UNIQUE-UUID" // Fallback placeholder
	}
	externalProvider := provider.NewWebhookProvider(webhookURL)

	go func() {
		for {
			InQueue, err := rdb.GetAndRemoveReadyTasks(ctx, "notification_retry_set", time.Now().Unix())
			if err == nil && len(InQueue) > 0 {
				for _, id := range InQueue {
					rdb.PushToQueue(ctx, "notification_queue_high", id)
				}
				log.Printf("⏰ Janitor: Moved %d messages from delay to main queue", len(InQueue))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	log.Println("👷 Worker started, listening for notifications...")

	for {
		// Define queue priority order
		queues := []string{
			"notification_queue_high",
			"notification_queue_normal",
			"notification_queue_low",
		}

		// BLPop prioritizes queues from left to right
		result, err := rdb.PopFromQueues(ctx, queues...)
		if err != nil {
			continue
		}

		notificationID := result[1]
		log.Printf("Processing %s: %s", result[0], notificationID)

		// 1. Fetch notification details
		n, err := repo.GetByID(ctx, notificationID)
		if err != nil {
			log.Printf("Failed to fetch notification %s: %v", notificationID, err)
			continue
		}

		if n.Status == domain.StatusCancelled {
			log.Printf("❌ Notification %s is cancelled, skipping provider send...", notificationID)
			continue
		}

		// Check rate limit (100 msgs/sec per channel)
		isLimited, err := rdb.IsRateLimited(ctx, string(n.Channel), 100)
		if err != nil {
			log.Printf("Rate limit check error: %v", err)
			continue
		}

		if isLimited {
			// If limited, push the ID back to the queue and wait a bit
			log.Printf("⚠️ Rate limit reached for %s, re-queuing %s", n.Channel, n.ID)
			queueName := fmt.Sprintf("notification_queue_%s", n.Priority)
			rdb.PushToQueue(ctx, queueName, n.ID.String())

			time.Sleep(100 * time.Millisecond) // Cool down
			continue
		}

		// 2. Send via external provider
		res, err := externalProvider.Send(ctx, n)
		if err != nil {
			log.Printf("❌ Delivery failed for %s: %v", n.ID, err)

			if n.RetryCount < 3 {
				backoffSchedule := []int{1, 5, 15}
				delay := time.Duration(backoffSchedule[n.RetryCount]) * time.Minute

				nextRetry := time.Now().Add(delay)
				newRetryCount := n.RetryCount + 1

				repo.UpdateRetryInfo(ctx, n.ID, newRetryCount, nextRetry)

				rdb.AddToDelayedQueue(ctx, "notification_retry_set", n.ID.String(), nextRetry.Unix())
				log.Printf("🔄 Moved to delayed queue for %v", delay)
			} else {
				repo.UpdateStatus(ctx, n.ID.String(), domain.StatusFailed)
				log.Printf("🚫 Max retries (3) reached for %s. Marked as FAILED.", n.ID)
			}
			continue
		}

		log.Printf("🚀 Provider Accepted: %s (ExternalID: %s)", notificationID, res.ExternalID)

		// 3. Update Status in DB
		err = repo.UpdateStatus(ctx, notificationID, domain.StatusSent)
		if err != nil {
			log.Printf("Failed to update status for %s: %v", notificationID, err)
			continue
		}

		log.Printf("✅ Notification %s marked as SENT", notificationID)
	}
}
