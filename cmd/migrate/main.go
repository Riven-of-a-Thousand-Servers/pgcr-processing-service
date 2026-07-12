package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"pgcr-processing-service/internal/db"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	_, err := db.Connect(ctx)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
}
