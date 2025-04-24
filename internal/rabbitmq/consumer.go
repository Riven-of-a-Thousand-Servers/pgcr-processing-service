package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"pgcr-processing-service/internal/mapper"
	"pgcr-processing-service/internal/service"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	amqp091 "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Channel     *amqp091.Channel
	Msgs        <-chan amqp091.Delivery
	PgcrService service.PgcrService
	PgcrMapper  mapper.PgcrMapper
}

func NewConsumer(url, queueName string, channel *amqp091.Channel) (*Consumer, error) {
	msgs, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("Unable to declare a queue with name [%s]: %v", queueName, err)
	}

	return &Consumer{
		Channel: channel,
		Msgs:    msgs,
	}, nil
}

// Consume messages from a rabbitmq queue and process PGCRs into its corresponding database tables
func (c *Consumer) Consume() error {
	for delivery := range c.Msgs {
		var pgcr types.PostGameCarnageReportResponse
		err := json.Unmarshal(delivery.Body, &pgcr)
		if err != nil {
			return fmt.Errorf("Error unmarshalling body from message: %v", err)
		}

		instanceId := pgcr.Response.ActivityDetails.InstanceId
		log.Printf("Processing pgcr [%s]...", instanceId)
		_, ppgcr, err := c.PgcrMapper.ToProcessedPgcr(&pgcr.Response)
		if err != nil {
			return fmt.Errorf("Error mapping pgcr [%s] to a processed pgcr: %v", instanceId, err)
		}
		err = c.PgcrService.ProcessPgcr(*ppgcr)
		if err != nil {
			return fmt.Errorf("Error processing ppgcr [%d] into database tables: %v", ppgcr.InstanceId, err)
		}
		log.Printf("Finished processing pgcr [%s]!", instanceId)
	}
	return nil
}
