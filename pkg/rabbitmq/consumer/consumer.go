package consumer

import (
	"coffee-shop/pkg/custom_errors"
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_exchangeKind       = "direct"
	_exchangeDurable    = true
	_exchangeAutoDelete = false
	_exchangeInternal   = false
	_exchangeNoWait     = false

	_queueDurable    = true
	_queueAutoDelete = false
	_queueExclusive  = false
	_queueNoWait     = false

	_prefetchCount  = 5
	_prefetchSize   = 0
	_prefetchGlobal = false

	_consumeAutoAck   = false
	_consumeExclusive = false
	_consumeNoLocal   = false
	_consumeNoWait    = false

	_exchangeName   = "orders-exchange"
	_queueName      = "orders-queue"
	_bindingKey     = "orders-routing-key"
	_consumerTag    = "orders-consumer"
	_workerPoolSize = 24
)

type consumer struct {
	exchangeName, queueName, bindingKey, consumerTag string
	workerPoolSize                                   int
	amqpConn                                         *amqp.Connection
}

func NewConsumer(amqConn *amqp.Connection) EventConsumer {
	return &consumer{
		amqpConn:       amqConn,
		exchangeName:   _exchangeName,
		queueName:      _queueName,
		bindingKey:     _bindingKey,
		consumerTag:    _consumerTag,
		workerPoolSize: _workerPoolSize,
	}
}

func (c *consumer) Configure(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *consumer) StartConsumer(worker Worker) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.createChannel()
	if err != nil {
		return err
	}

	defer ch.Close()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		_consumeAutoAck,
		_consumeExclusive,
		_consumeNoLocal,
		_consumeNoWait,
		nil,
	)

	if err != nil {
		return custom_errors.InternalError.Wrap(nil, false, err, "Error while consuming messages")
	}
	forever := make(chan bool)

	for i := 0; i < c.workerPoolSize; i++ {
		go worker.Start(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp.Error))
	log.Print("ch.NotifyClose", chanErr)
	<-forever

	return chanErr

}

func (c *consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(nil, false, err, "Error while creating channel")
	}

	err = ch.ExchangeDeclare(c.exchangeName,
		_exchangeKind,
		_exchangeDurable,
		_exchangeAutoDelete,
		_exchangeInternal,
		_exchangeNoWait,
		nil)
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(nil, false, err, "Error while declaring exchange")
	}

	queue, err := ch.QueueDeclare(
		c.queueName,
		_queueDurable,
		_queueAutoDelete,
		_queueExclusive,
		_queueNoWait,
		nil)

	if err != nil {
		return nil, custom_errors.InternalError.Wrap(nil, false, err, "Error while declaring queue")
	}

	err = ch.QueueBind(queue.Name, c.bindingKey, c.exchangeName, _queueNoWait, nil)

	if err != nil {
		return nil, custom_errors.InternalError.Wrap(nil, false, err, "Error while binding queue")
	}

	err = ch.Qos(_prefetchCount, _prefetchSize, _prefetchGlobal)

	if err != nil {
		return nil, custom_errors.InternalError.Wrap(nil, false, err, "Error while setting QOS")
	}

	return ch, nil

}
