package configurations

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/contracts"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	"github.com/vmdt/notification-worker/server/consumer"
	"github.com/vmdt/notification-worker/shared"
)

func ConfigConsumers(
	ctx context.Context,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config,
	mailer *mailer.Mailer,
	publisher rabbitmq.IPublisher,
	echo *echo.Echo,
	a *asynq.Client,
	notificationScheduleRepository contracts.NotificationScheduleRepository,

) error {
	discountBase := shared.DiscountBase{
		Log:                            log,
		Cfg:                            cfg,
		ConnRabbitmq:                   connRabbitmq,
		Publisher:                      publisher,
		NotificationScheduleRepository: notificationScheduleRepository,
		Ctx:                            ctx,
		Mailer:                         mailer,
		Echo:                           echo,
	}

	microserviceBase := shared.MicroserviceBase{
		Log:         log,
		Cfg:         cfg,
		AsynqClient: a,
	}

	microServiceConsumer := rabbitmq.NewConsumer(ctx, cfg.Rabbitmq, connRabbitmq, log, consumer.HandleConsumeMicroserviceMessage)

	discountConsumer := rabbitmq.NewConsumer(ctx, cfg.Rabbitmq, connRabbitmq, log, consumer.HandleConsumeDiscountMessage)

	go func() {
		err := discountConsumer.ConsumeMessage(nil, "discount", "discount_queue", "discount_key", &discountBase)
		if err != nil {
			log.Error(err)
		}
	}()

	go func() {
		err := microServiceConsumer.ConsumeMessage(nil, cfg.Rabbitmq.ExchangeName, "microservice_queue", "microservice_key", &microserviceBase)
		if err != nil {
			log.Error(err)
		}
	}()

	log.Info("RabbitMQ consumer started")
	return nil
}
