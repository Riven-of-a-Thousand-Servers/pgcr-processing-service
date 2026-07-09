package rabbitmq

import (
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	Queue   amqp091.Queue
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

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare queue", "Error", err)
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
		Queue:   q,
	}, err
}

// Instantiate a queue Consumer
// The name parameter declares the name of the consumer
func (r *RabbitMQ) Consumer(name string) (<-chan amqp091.Delivery, error) {
	delivery, err := r.Channel.Consume(r.Queue.Name, name, false, false, false, false, nil)
	if err != nil {
		slog.Error("Error declaring consumer for RabbitMQ", "Error", err)
		return nil, err
	}

	return delivery, nil
}
