package configurations

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	"github.com/vmdt/notification-worker/server/consumer"
)

func ConfigConsumers(
	ctx context.Context,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config,
) error {
	discountConsumer := rabbitmq.NewConsumer(ctx, cfg.Rabbitmq, connRabbitmq, log, consumer.HandleConsumeDiscountMessage)

	go func() {
		err := discountConsumer.ConsumeMessage(nil, cfg.Rabbitmq.ExchangeName, "discount_queue", "discount_key")
		if err != nil {
			log.Error(err)
		}
	}()

	log.Info("RabbitMQ consumer started")
	return nil
}
