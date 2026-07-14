package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"pgcr-processing-service/internal/db"
	"pgcr-processing-service/internal/processing"
	"pgcr-processing-service/internal/rabbitmq"
	"pgcr-processing-service/internal/redis"
)

var (
	rabbitMQUrl = "amqp://rabbitmq:5672"
	postgresUrl = "postgres://%s:%s@postgres:5432/postgres?sslmode=disable"
	redisUrl    = "redis:6379"
	goroutines  = 100
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()
	rabbitmq, err := rabbitmq.Connect("rivenbot_pgcr", rabbitMQUrl)
	if err != nil {
		slog.Error("Error happened while connecting to RabbitMQ", "error", err)
		os.Exit(1)
	}

	conn, err := db.Connect(ctx, postgresUrl)
	if err != nil {
		slog.Error("Error happened while connecting to DB", "error", err)
		os.Exit(1)
	}

	queries, err := db.Prepare(ctx, conn)
	if err != nil {
		slog.Error("Error creating and preparing queries", "error", err)
		os.Exit(1)
	}

	redis := redis.NewService(redisUrl)
	processor := processing.NewPgcrProcessor(conn, queries, rabbitmq, redis)

	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			slog.Info("Starting worker", "Id", i)
			_ = processor.StartWork(ctx, i)
			slog.Info("Shutting down worker", "Id", i)
		})
	}

	wg.Wait()
}
