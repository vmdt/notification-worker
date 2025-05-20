package consumer

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/shared"
)

type DiscountMessage struct {
	To       string                 `json:"to"`
	Subject  string                 `json:"subject"`
	From     string                 `json:"from"`
	Template string                 `json:"template"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func HandleConsumeDiscountMessage(queueName string, msg amqp.Delivery, dependencies *shared.DiscountBase) error {
	log.Infof("Message received on queue: %s with message: %s", queueName, string(msg.Body))

	var discountMsg DiscountMessage
	if err := json.Unmarshal(msg.Body, &discountMsg); err != nil {
		log.Errorf("Error unmarshaling message: %v", err)
		return err
	}

	err := dependencies.Mailer.SendMail(
		discountMsg.Template,
		discountMsg.To,
		discountMsg.Metadata,
	)

	if err != nil {
		log.Errorf("Error sending email: %v", err)
		return err
	}

	return nil
}
