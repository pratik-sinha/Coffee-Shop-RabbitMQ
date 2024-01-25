package rabbitmq

import (
	"coffee-shop/internal/kitchen/events"
	"coffee-shop/internal/kitchen/service"
	"coffee-shop/pkg/rabbitmq/consumer"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type kitchenWorker struct {
	ks service.KitchenService
}

func NewKitchenWorker(k service.KitchenService) consumer.Worker {
	return &kitchenWorker{ks: k}
}

func (k *kitchenWorker) Start(ctx context.Context, messages <-chan amqp091.Delivery) {
	for delivery := range messages {
		log.Print("processDeliveries", "delivery_tag", delivery.DeliveryTag)
		log.Print("received", "delivery_type", delivery.Type)

		switch delivery.Type {
		case "kitchen-order-created":
			var payload events.ItemsOrderedEvent

			err := json.Unmarshal(delivery.Body, &payload)
			if err != nil {
				log.Print("failed to Unmarshal message", err)
			}

			err = k.ks.ProcessItems(ctx, payload)

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
