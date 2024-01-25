package rabbitmq

import (
	"errors"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2
)

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

func NewRabbitMQConn(url string) (*amqp.Connection, func(), error) {
	var counts int32

	for {
		amqConn, err := amqp.Dial(url)
		cleanUp := func() { amqConn.Close() }
		if err != nil {
			log.Printf("Failed to connect to rabbitmq server at: %s", url)
			counts++
		} else {
			return amqConn, cleanUp, nil
		}

		if counts > _retryTimes {
			return nil, nil, ErrCannotConnectRabbitMQ
		}
		log.Print("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)
	}

}
