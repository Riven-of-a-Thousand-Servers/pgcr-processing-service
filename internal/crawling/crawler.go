package crawling

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"pgcr-processing-service/internal/rabbitmq"

	"github.com/rabbitmq/amqp091-go"
)

var (
	keyHeaderName = "x-api-key"
	baseUrl       = "https://stats.bungie.net/Platform/Destiny2/Stats/PostGameCarnageReport/%d/"
)

type PgcrCrawler struct {
	Offset   int64
	MaxSize  int64
	In       <-chan int64
	Client   *http.Client
	Rabbitmq *rabbitmq.RabbitMQ
}

func NewPgcrCrawler(rabbitmq *rabbitmq.RabbitMQ, client *http.Client, gen <-chan int64, maxSize int64) *PgcrCrawler {
	return &PgcrCrawler{
		Rabbitmq: rabbitmq,
		Client:   client,
		In:       gen,
		MaxSize:  maxSize,
	}
}

func (c *PgcrCrawler) Crawl(ctx context.Context, id int64, apiKey string) {
	ch, err := c.Rabbitmq.Conn.Channel()
	if err != nil {
		slog.Error("Failed to oppen rabbitmq channel", "workerId", id, "error", err)
	}
	defer ch.Close()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Crawler instance shutting down", "workerId", id)
			return
		case next, ok := <-c.In:
			if !ok {
				slog.Info("Input channel closed. Exiting", "workerId", id)
				return
			}
			slog.Info("Worker processing pgcr", "workerId", id, "pgcr", next)
			url := fmt.Sprintf(baseUrl, next)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				slog.Error("Worker unable to create requests. Exiting.", "workerId", id, "error", err)
				return
			}

			req.Header.Add(keyHeaderName, apiKey)
			res, err := c.Client.Do(req)
			if err != nil {
				slog.Error("Unable to get a response from Bungie. Exiting.", "workerId", id, "error", err)
				return
			}

			// Stop if there's errors reading HTTP bodies from requests
			data, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				slog.Error("Error reading response body", "error", err)
				return
			}

			slog.Debug("Response raw data", "data", string(data))
			if int64(len(data)) > c.MaxSize {
				slog.Error(fmt.Sprintf("Response exceeded limit of %d bytes, refusing to process", c.MaxSize), "workerId", id, "pgcr", next)
				continue
			}

			publishing := amqp091.Publishing{
				MessageId: strconv.FormatInt(next, 10),
				Headers: map[string]any{
					"source": "Crawler",
				},
				ContentType:     "application/json",
				ContentEncoding: "utf-8",
				Body:            data,
			}

			if err := ch.PublishWithContext(ctx, "", c.Rabbitmq.Queue.Name, false, false, publishing); err != nil {
				slog.Error("Unable to publish message", "messageId", next, "crawlerId", id)
			}

			slog.Info("Successfully published pgcr", "pgcr", next)
		}
	}
}
