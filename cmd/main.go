package main

import (
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/contracts/repositories"
	"github.com/vmdt/notification-worker/pkg/cron"
	echoserver "github.com/vmdt/notification-worker/pkg/echo"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/http"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/mongodb"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	"github.com/vmdt/notification-worker/server"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				cron.NewCronManager,
				http.NewContext,
				echoserver.NewEchoServer,
				mongodb.NewMongoDB,
				repositories.NewMongoNotificationScheduleRepository,
				mailer.NewMailer,
			),
			fx.Invoke(server.RunServers),
		),
	).Run()
}
