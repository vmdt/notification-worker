package consumer

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func HandleConsumeDiscountMessage(queueName string, msg amqp.Delivery) error {
	log.Infof("Message received on queue: %s with message: %s", queueName, string(msg.Body))

	return nil
}
