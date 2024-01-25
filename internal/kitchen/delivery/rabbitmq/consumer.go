package rabbitmq

import (
	"coffee-shop/internal/kitchen/service"
	"coffee-shop/pkg/rabbitmq/consumer"
)

type KitchenConsumer interface {
	Configure()
	StartConsumer(service.KitchenService) error
}

type kitchenConsumer struct {
	con consumer.EventConsumer
}

func NewKitchenConsumer(con consumer.EventConsumer) KitchenConsumer {
	return &kitchenConsumer{con: con}
}

func (k *kitchenConsumer) Configure() {
	k.con.Configure(consumer.ExchangeName("kitchen-order-exchange"),
		consumer.QueueName("kitchen-order-queue"),
		consumer.BindingKey("kitchen-order-routing-key"),
		consumer.ConsumerTag("kitchen-order-consumer"))
}

func (k *kitchenConsumer) StartConsumer(ks service.KitchenService) error {
	worker := NewKitchenWorker(ks)
	err := k.con.StartConsumer(worker)
	if err != nil {
		return err
	}
	return nil
}
