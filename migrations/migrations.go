package migrations

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

// RunMigrations executes all pending database migrations
func RunMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	log.Println("checking database migrations...")
	if err := goose.Up(db, "."); err != nil {
		return err
	}

	log.Println("migrations completed successfully")
	return nil
}
