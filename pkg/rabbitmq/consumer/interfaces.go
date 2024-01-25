package consumer

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Worker interface {
	Start(ctx context.Context, messages <-chan amqp.Delivery)
}

type EventConsumer interface {
	Configure(...Option)
	StartConsumer(Worker) error
}
