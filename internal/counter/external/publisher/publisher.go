package external

import (
	"coffee-shop/pkg/rabbitmq/publisher"
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type (
	baristaEventPublisher struct {
		pub publisher.EventPublisher
	}
	kitchenEventPublisher struct {
		pub publisher.EventPublisher
	}
)

func NewBaristaEventPublisher(amq *amqp091.Connection) (BaristaEventPublisher, error) {
	pub, err := publisher.NewPublisher(amq)
	if err != nil {
		return nil, err
	}
	return &baristaEventPublisher{
		pub: pub,
	}, nil
}

func (b *baristaEventPublisher) Configure() {
	b.pub.Configure(
		publisher.ExchangeName("barista-order-exchange"),
		publisher.BindingKey("barista-order-routing-key"),
		publisher.MessageTypeName("barista-order-created"),
	)
}

func (b *baristaEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return b.pub.Publish(ctx, body, contentType)
}

func NewKitchenEventPublisher(amq *amqp091.Connection) (KitchenEventPublisher, error) {
	pub, err := publisher.NewPublisher(amq)
	if err != nil {
		return nil, err
	}
	return &kitchenEventPublisher{
		pub: pub,
	}, nil
}

func (k *kitchenEventPublisher) Configure() {
	k.pub.Configure(
		publisher.ExchangeName("kitchen-order-exchange"),
		publisher.BindingKey("kitchen-order-routing-key"),
		publisher.MessageTypeName("kitchen-order-created"),
	)
}

func (k *kitchenEventPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	return k.pub.Publish(ctx, body, contentType)
}
