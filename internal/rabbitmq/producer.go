package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
)

type QueueProducer interface {
	SendVerifyEmailTask(ctx context.Context, payload VerifyEmailPayload) error
}

type Producer struct {
	Ch *amqp.Channel
}

func NewProducer(url string) (QueueProducer, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		VerifyEmailKey, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		return nil, err
	}

	err = ch.Qos(1, 0, false)

	if err != nil {
		return nil, err
	}

	return Producer{
		Ch: ch,
	}, err
}
