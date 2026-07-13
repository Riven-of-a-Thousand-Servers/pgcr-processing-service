package rabbitmq

import (
	"context"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn  *amqp091.Connection
	Queue amqp091.Queue
}

func Connect(queueName, url string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		slog.Error("Error dialing RabbitMQ", "Error", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		slog.Error("Failed to open a concurrent channel", "Error", err)
		conn.Close()
		return nil, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare queue", "Error", err)
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn:  conn,
		Queue: q,
	}, err
}

// Instantiate a queue Consumer
// The name parameter declares the name of the consumer
func (r *RabbitMQ) Consumer(ctx context.Context, consumerName string) (<-chan amqp091.Delivery, *amqp091.Channel, error) {
	ch, err := r.Conn.Channel()
	if err != nil {
		slog.Error("Failed to open amqp channel", "error", err, "consumer", consumerName)
		return nil, nil, err
	}

	delivery, err := ch.ConsumeWithContext(ctx, r.Queue.Name, consumerName, false, false, false, false, nil)
	if err != nil {
		slog.Error("Error declaring consumer for RabbitMQ", "Error", err)
		return nil, nil, err
	}

	return delivery, ch, nil
}
