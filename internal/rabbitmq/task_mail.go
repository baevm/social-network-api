package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type VerifyEmailPayload struct {
	Username string
	Email    string
}

const VerifyEmailKey = "verify_email"

func (p Producer) SendVerifyEmailTask(ctx context.Context, payload VerifyEmailPayload) error {
	body, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	err = p.Ch.Publish(
		"",
		VerifyEmailKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})

	return err
}

func (c Consumer) ConsumeVerifyEmailTask(msgs <-chan amqp.Delivery) error {
	for d := range msgs {
		var payload VerifyEmailPayload

		c.logger.Infoln("Received a message: ", string(d.Body), " from: ", d.RoutingKey)

		err := json.Unmarshal(d.Body, &payload)

		if err != nil {
			continue
		}

		//verifyEmail, err := c.userRepo.CreateVerifyEmail(ctx, payload.Username, payload.Email)

		c.logger.Infoln("Sending email to: ", payload.Username)

		verifyLink := fmt.Sprintf("http://localhost:5000/v1/verify_email?email_id=%d&secret_code=%s", 1, "123")

		subject := "Confirm your email"
		content := fmt.Sprintf(`
		<h1>Confirm your email by clicking this link</h1>
		<div>
		<a href="%s">Click</a>
		</div>
		`, verifyLink)

		err = c.mailer.SendEmail(subject, content, []string{payload.Email}, nil, nil, nil)

		if err != nil {
			c.logger.Errorln("Error sending email: ", err)
		}
	}

	return nil
}
