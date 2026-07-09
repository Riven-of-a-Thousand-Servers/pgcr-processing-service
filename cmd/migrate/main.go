package main

import (
	"context"
	"embed"
	"log/slog"
	"os/signal"
	"pgcr-processing-service/internal/db"
	"syscall"

	"sync"

	"github.com/pressly/goose/v3"
)

//go:embed: ../../db/migrations/*.sql
var migrations embed.FS

var gooseInitErr error
var gooseOnce sync.Once

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	db, err := db.Connect(ctx)
	if err != nil {
		slog.Error("Error connecting to database", "Error", err)
	}
	defer db.Close()

	if err := goose.Up(db, "db/migrations"); err != nil {
		db.Close()
		slog.Error("Error running migrations with Goose", "Error", err)
	}
}

func initGoose() error {
	gooseOnce.Do(func() {
		goose.SetBaseFS(migrations)
		gooseInitErr = goose.SetDialect("pgx")
	})
	return gooseInitErr
}
