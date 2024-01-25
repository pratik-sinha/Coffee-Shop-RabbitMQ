package publisher

import (
	"coffee-shop/pkg/custom_errors"
	"context"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "orders-exchange"
	_bindingKey      = "orders-routing-key"
	_messageTypeName = "ordered"
)

type publisher struct {
	exchangeName, bindingKey string
	messageTypeName          string
	amqpConn                 *amqp.Connection
}

func NewPublisher(amqpConn *amqp.Connection) (EventPublisher, error) {
	pub := &publisher{
		amqpConn:        amqpConn,
		exchangeName:    _exchangeName,
		bindingKey:      _bindingKey,
		messageTypeName: _messageTypeName,
	}

	return pub, nil
}

func (p *publisher) Configure(opts ...Option) {
	for _, opt := range opts {
		opt(p)
	}
}

func (p *publisher) Publish(ctx context.Context, body []byte, contentType string) error {
	ch, err := p.amqpConn.Channel()
	if err != nil {
		return custom_errors.InternalError.Wrap(nil, false, err, "Error while creating channel")
	}
	defer ch.Close()

	if err = ch.PublishWithContext(ctx, p.exchangeName, p.
		bindingKey, _publishMandatory, _publishImmediate, amqp.Publishing{
		ContentType:  contentType,
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.NewString(),
		Body:         body,
		Type:         p.messageTypeName,
	}); err != nil {
		return custom_errors.InternalError.Wrap(nil, false, err, "Error while publishing message")
	}

	return nil
}
