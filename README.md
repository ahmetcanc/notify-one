# Notify One 🚀

High-performance, scalable notification service built with Go.
Designed to handle single and batch notifications with a focus on reliability, exactly-once delivery, and priority-based processing.

---

## 🏗 Architecture Decisions
Clean Architecture: The project is organized into Domain, Usecase, API, and Infrastructure layers to ensure separation of concerns and high testability.

Janitor Pattern (Reliability): Failed notifications (due to provider rate limits or temporary outages) are moved to a Redis Sorted Set (ZSET) with a visibility timeout.

Exponential Backoff: Retries are scheduled with increasing delays (e.g., 1m, 5m, 15m) to prevent overwhelming third-party providers.

Idempotency: Supports idempotency_key to prevent duplicate notification delivery in case of network retries.

---

# 🛠 Tech Stack

* **Language:** Go 1.25+
* **Database:** PostgreSQL 16 (`pgx` pool)
* **Cache:** Redis 7
* **Migrations:** Goose
* **Infrastructure:** Docker & Docker Compose

---

# 🚀 Quick Start

## 1️⃣ Environment Setup

Create a `.env` file in the project root:
`cp .env.example .env`

```bash
APP_PORT=3333

DB_HOST=postgres
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=notifydb
DB_SSLMODE=disable

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 📚 API Documentation

```bash
👉 http://localhost:3333/swagger/index.html
```

---

## TESTS
```bash
go test ./... -v
```

---

## 2️⃣ Launch the Project

Build and start all services (API, Worker, DB, Redis):
```bash
docker compose up --build
```

---

## 3️⃣ Management & Cleanup

Stop services:
```bash
docker compose down
```

Full reset (including volumes):
```bash
docker compose down -v && docker compose up --build
```

---

## 4️⃣ Send Notifications

```bash
# Single Notification (with Idempotency Key)
curl -X POST http://localhost:3333/api/v1/notifications \
     -H "Content-Type: application/json" \
     -d '{
       "recipient": "example@example.com",
       "channel": "email",
       "content": "Hello World!",
       "priority": "high",
       "idempotency_key": "unique-key-101"
     }'

# Batch Notifications
curl -X POST http://localhost:3333/api/v1/notifications/batch \
     -H "Content-Type: application/json" \
     -d '[
       {"recipient": "user1@test.com", "channel": "sms", "content": "Batch 1", "priority": "normal"},
       {"recipient": "user2@test.com", "channel": "sms", "content": "Batch 2", "priority": "normal"}
     ]'
```

---

## 5️⃣ Management & Filtering

```bash
# List all SENT notifications
curl "http://localhost:3333/api/v1/notifications?status=sent&limit=10"

# Filter by Channel and Priority
curl "http://localhost:3333/api/v1/notifications?channel=sms&priority=high"

# Bulk Cancel a Batch (Replace {id} with the batch_id from the POST response)
curl -X PATCH http://localhost:3333/api/v1/notifications/batch/{id}/cancel
```

---

## 🗄️ Database Migrations

We use Goose for versioned database schema changes.

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

goose -dir migrations create <migration_name> sql
```

---

# 🔌 API Endpoints

| Endpoint                                  | Method | Description                                          | Status     |
| ------------------------------------------| ------ | -----------------------------------------------------| ---------- |
| `/health`                                 | GET    | System health & connectivity check                   | ✅ Ready   |
| `/api/v1/metrics`                         | GET    | Real-time queue depths and system metrics            | ✅ Ready   |
| `/api/v1/notifications`                   | POST   | Send single notification with Idempotency support    | ✅ Ready   |
| `/api/v1/notifications/batch`             | POST   | Send batch notifications with async processing       | ✅ Ready   |
| `/api/v1/notifications`                   | GET    | List and filter notifications with pagination        | ✅ Ready   |
| `/api/v1/notifications/batch/{id}/cancel` | PATCH  | Cancel pending notifications within a specific batch | ✅ Ready   |

---

# 📂 Project Structure

```
.
├── cmd/
│   ├── api/              # API entry point (main.go)
│   └── worker/           # Background worker entry point
├── internal/
    ├── api               # HTTP Handlers (Encapsulated)
│   ├── domain/           # Business logic & entities
│   ├── infrastructure/   # DB, cache, config, provider implementations
│   └── usecase/          # Application orchestration (Business rules)
├── migrations/           # SQL migration files & runner
├── docker-compose.yml    # Infrastructure orchestration
└── Dockerfile            # Multi-stage build for Go
├── .env                  # Environment variables
└── Readme.md             # Project documentation
```

---

# 📜 Database Schema

The system uses custom **PostgreSQL ENUM types** for strict validation.

### Channels

* `sms`
* `email`
* `push`

### Status

* `pending`
* `processing`
* `sent`
* `failed`
* `cancelled`

### Priority

* `low`
* `normal`
* `high`

---

### Useful Commands
To access the postgreSQL database via terminal:
```bash
docker exec -it notify_postgres psql -U user -d notifydb

```
To access the Redis database via terminal:
```bash
docker exec -it notify_redis redis-cli
```

---

## 👤 Author
Ahmet Can Ceylan - [canceylan.dev](https://canceylan.dev)