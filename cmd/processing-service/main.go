package pgcrprocessingservice

import (
	"context"
	"log/slog"

	"pgcr-processing-service/internal/db"
	"pgcr-processing-service/internal/rabbitmq"
)

var rabbitMQUrl = "amqp://user:password@localhost:5672"

func main() {
	ctx := context.Background()
	_, err := rabbitmq.Connect("rivenbot_pgcr", rabbitMQUrl)
	if err != nil {
		slog.Error("Error happened while connecting to RabbitMQ", "Error", err)
	}

	_, err = db.Connect(ctx)
	if err != nil {
		slog.Error("Error happened while connecting to DB", "Error", err)
	}
}
