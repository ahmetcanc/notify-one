# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-bin cmd/api/main.go
RUN go build -o worker-bin cmd/worker/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-bin .
COPY --from=builder /app/worker-bin .
COPY --from=builder /app/migrations ./migrations
COPY .env .

EXPOSE 3333