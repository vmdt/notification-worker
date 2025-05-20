package shared

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/contracts"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
)

type DiscountBase struct {
	Log                            logger.ILogger
	Cfg                            *config.Config
	Publisher                      rabbitmq.IPublisher
	ConnRabbitmq                   *amqp.Connection
	Echo                           *echo.Echo
	Mailer                         *mailer.Mailer
	NotificationScheduleRepository contracts.NotificationScheduleRepository
	Ctx                            context.Context
}
