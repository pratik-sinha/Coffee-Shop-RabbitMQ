package rabbitmq

import (
	"coffee-shop/internal/counter/service"
	"coffee-shop/pkg/rabbitmq/consumer"
)

type CounterConsumer interface {
	Configure()
	StartConsumer(service.CounterService) error
}

type counterConsumer struct {
	con consumer.EventConsumer
}

func NewCounterConsumer(con consumer.EventConsumer) CounterConsumer {
	return &counterConsumer{con: con}
}

func (c *counterConsumer) Configure() {
	c.con.Configure(consumer.ExchangeName("counter-order-exchange"),
		consumer.QueueName("counter-order-queue"),
		consumer.BindingKey("counter-order-routing-key"),
		consumer.ConsumerTag("counter-order-consumer"))
}

func (c *counterConsumer) StartConsumer(cs service.CounterService) error {
	worker := NewCounterWorker(cs)
	err := c.con.StartConsumer(worker)
	if err != nil {
		return err
	}
	return nil
}
