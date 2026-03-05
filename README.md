# Notify One 🚀

High-performance, scalable notification service built with **Go**.
Designed to handle **single and batch notifications** using **PostgreSQL** and **Redis**.

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

# 🔌 API Endpoints

| Endpoint                      | Method | Description                        | Status     |
| ----------------------------- | ------ | ---------------------------------- | ---------- |
| `/health`                     | GET    | System health & connectivity check | ✅ Ready   |
| `/api/v1/notifications`       | POST   | Send single notification (async)   | ✅ Ready   |
| `/api/v1/notifications/batch` | POST   | Send batch notifications (async)   | ✅ Ready   |

---

# 📂 Project Structure

```
.
├── cmd/
│   ├── api/              # API entry point (main.go)
│   └── worker/           # Background worker entry point
├── internal/
│   ├── domain/           # Business logic & entities
│   ├── infrastructure/   # DB, cache, config, provider implementations
│   └── usecase/          # Application orchestration (Business rules)
├── migrations/           # SQL migration files & runner
├── pkg/                  # Shared libraries
│   ├── logger/           # Structured logging implementation
│   └── metrics/          # Prometheus or custom metrics
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

### Useful Commands
To access the postgreSQL database via terminal:
```bash
docker exec -it notify_postgres psql -U user -d notifydb

```
To access the Redis database via terminal:
```bash
docker exec -it notify_redis redis-cli
