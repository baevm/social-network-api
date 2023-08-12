package rabbitmq

import (
	"social-network-api/internal/mail"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type Consumer struct {
	Ch     *amqp.Channel
	logger *zap.SugaredLogger
	mailer mail.EmailSender
}

func NewConsumer(url string, logger *zap.SugaredLogger, mailer mail.EmailSender) (*Consumer, error) {
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

	return &Consumer{
		Ch:     ch,
		logger: logger,
		mailer: mailer,
	}, nil
}

func (c Consumer) Start() error {
	var stop chan struct{}

	msgs, err := c.Ch.Consume(
		VerifyEmailKey,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	c.ConsumeVerifyEmailTask(msgs)

	<-stop

	return err
}
