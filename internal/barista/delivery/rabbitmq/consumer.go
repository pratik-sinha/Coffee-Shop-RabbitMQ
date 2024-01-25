package rabbitmq

import (
	"coffee-shop/internal/barista/service"
	"coffee-shop/pkg/rabbitmq/consumer"
)

type BaristaConsumer interface {
	Configure()
	StartConsumer(service.BaristaService) error
}

type baristaConsumer struct {
	con consumer.EventConsumer
}

func NewBaristaConsumer(con consumer.EventConsumer) BaristaConsumer {
	return &baristaConsumer{con: con}
}

func (b *baristaConsumer) Configure() {
	b.con.Configure(consumer.ExchangeName("barista-order-exchange"),
		consumer.QueueName("barista-order-queue"),
		consumer.BindingKey("barista-order-routing-key"),
		consumer.ConsumerTag("barista-order-consumer"))
}

func (b *baristaConsumer) StartConsumer(bs service.BaristaService) error {
	worker := NewBaristaWorker(bs)
	err := b.con.StartConsumer(worker)
	if err != nil {
		return err
	}
	return nil
}
