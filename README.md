# Notify One 🚀

High-performance, scalable notification service built with **Go**.
Designed to handle **single and batch notifications** using **PostgreSQL** and **Redis**.

---

# 🛠 Tech Stack

* **Language:** Go 1.23+
* **Database:** PostgreSQL 15 (`pgx` pool)
* **Cache:** Redis 7
* **Migrations:** Goose
* **Infrastructure:** Docker & Docker Compose

---

# 🚀 Quick Start

## 1️⃣ Environment Setup

Create a `.env` file in the project root:

```bash
APP_PORT=3333

DB_HOST=localhost
DB_PORT=5555
DB_USER=user
DB_PASSWORD=password
DB_NAME=notifydb
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

---

## 2️⃣ Start Infrastructure

Launch PostgreSQL and Redis containers:

```bash
docker compose up -d
```

---

## 3️⃣ Run the Application

The application automatically runs database migrations on startup using **embed.FS**.

```bash
go run cmd/api/main.go
```

---

# 🔌 API Endpoints

| Endpoint                      | Method | Description                        | Status     |
| ----------------------------- | ------ | ---------------------------------- | ---------- |
| `/health`                     | GET    | System health & connectivity check | ✅ Ready    |
| `/api/v1/notifications`       | POST   | Send single notification           | 🏗 Pending |
| `/api/v1/notifications/batch` | POST   | Send batch notifications           | 🏗 Pending |

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

