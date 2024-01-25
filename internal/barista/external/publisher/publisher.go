package external

import (
	"coffee-shop/pkg/rabbitmq/publisher"
	"context"
)

type counterEventPublisher struct {
	pub publisher.EventPublisher
}

func NewCounterEventPublisher(pub publisher.EventPublisher) CounterEventPublisher {
	return &counterEventPublisher{
		pub: pub,
	}
}

func (c *counterEventPublisher) Configure() {
	c.pub.Configure(
		publisher.ExchangeName("counter-order-exchange"),
		publisher.BindingKey("counter-order-routing-key"),
		publisher.MessageTypeName("barista-order-updated"),
	)
}

func (c *counterEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return c.pub.Publish(ctx, body, contentType)
}
