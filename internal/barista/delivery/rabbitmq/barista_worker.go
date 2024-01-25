package rabbitmq

import (
	"coffee-shop/internal/barista/events"
	"coffee-shop/internal/barista/service"
	"coffee-shop/pkg/rabbitmq/consumer"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type baristaWorker struct {
	bs service.BaristaService
}

func NewBaristaWorker(k service.BaristaService) consumer.Worker {
	return &baristaWorker{bs: k}
}

func (b *baristaWorker) Start(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		log.Print("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		log.Print("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "barista-order-created":
			var payload events.ItemsOrderedEvent

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				log.Print("failed to Unmarshal message", err)
			}

			err = b.bs.ProcessItems(ctx, payload)

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
