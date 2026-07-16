package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"pgcr-processing-service/internal/crawling"
	"pgcr-processing-service/internal/rabbitmq"
	"pgcr-processing-service/internal/transport"
	"pgcr-processing-service/internal/types/net"
	types "pgcr-processing-service/internal/types/rabbitmq"
)

var (
	goroutines = 50
	apiKey     = "27fd2725658c431992b1b0259682ad3c"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	rabbitmq, err := rabbitmq.Connect(types.RabbitQueueName, types.RabbitMQUrl)
	if err != nil {
		slog.Error("Error happened while connecting to RabbitMQ", "error", err)
		os.Exit(1)
	}
	defer rabbitmq.Conn.Close()

	client := http.Client{
		Transport: &transport.MaxSizeTransport{
			Base:    http.DefaultTransport,
			MaxSize: net.MAX_REQUEST_SIZE_KB,
		},
	}

	var wg sync.WaitGroup
	tick := time.NewTicker(10 * time.Second)

	wg.Add(1)
	in := make(chan int64, 100)
	go func(ctx context.Context, throttle *time.Ticker, start int64, in chan<- int64) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				slog.Info("Context cancelled. Exiting.")
				close(in)
				return
			case <-throttle.C:
				start += 1
				in <- int64(start)
			}
		}
	}(ctx, tick, 1, in)

	crawler := crawling.NewPgcrCrawler(rabbitmq, &client, in, net.MAX_REQUEST_SIZE_KB)
	for i := range goroutines {
		wg.Go(func() {
			crawler.Crawl(ctx, int64(i), apiKey)
		})
	}

	wg.Wait()
	slog.Info("All workers stopped, cleaning up resources")
}
