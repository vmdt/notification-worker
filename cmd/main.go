package main

import (
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/configurations"
	"github.com/vmdt/notification-worker/contracts/repositories"
	"github.com/vmdt/notification-worker/pkg/cron"
	echoserver "github.com/vmdt/notification-worker/pkg/echo"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/http"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/mongodb"
	"github.com/vmdt/notification-worker/pkg/queue"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	redis2 "github.com/vmdt/notification-worker/pkg/redis"
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
				redis2.NewRedisClient,
				queue.NewServeMux,
				queue.NewClient,
				queue.NewServer,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigConsumers),
			fx.Invoke(configurations.HookQueueClient),
			fx.Invoke(configurations.HookQueueServer),
		),
	).Run()
}
