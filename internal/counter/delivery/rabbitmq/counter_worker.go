package rabbitmq

import (
	"coffee-shop/internal/counter/events"
	"coffee-shop/internal/counter/service"
	"coffee-shop/pkg/rabbitmq/consumer"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type counterWorker struct {
	cs service.CounterService
}

func NewCounterWorker(cs service.CounterService) consumer.Worker {
	return &counterWorker{cs: cs}
}

func (c *counterWorker) Start(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		log.Print("processDeliveries", "delivery_tag ", delivery.DeliveryTag)
		log.Print("received", "delivery_type ", delivery.Type)
		log.Print("delivery_body", string(delivery.Body))

		switch delivery.Type {
		case "barista-order-updated":
			var payload events.ItemOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				log.Print("failed to Unmarshal message", err)
			}

			err = c.cs.UpdateOrder(ctx, payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					log.Print("failed to delivery.Reject", err)
				}

				log.Print("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					log.Print("failed to acknowledge delivery", err)
				}
			}
		case "kitchen-order-updated":
			var payload events.ItemOrderUpdated

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				log.Print("failed to Unmarshal message", err)
			}

			err = c.cs.UpdateOrder(ctx, payload)

			if err != nil {
				if err = delivery.Reject(false); err != nil {
					log.Print("failed to delivery.Reject", err)
				}

				log.Print("failed to process delivery", err)
			} else {
				err = delivery.Ack(false)
				if err != nil {
					log.Print("failed to acknowledge delivery", err)
				}
			}
		default:
			log.Print("default")
		}
	}

	log.Print("Deliveries channel closed")
}
