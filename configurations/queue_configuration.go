package configurations

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/contracts"
	mailer "github.com/vmdt/notification-worker/pkg/email"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	"github.com/vmdt/notification-worker/shared"
	"github.com/vmdt/notification-worker/workers"
	"go.uber.org/fx"
)

func HookQueueClient(lifecycle fx.Lifecycle, client *asynq.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return client.Close()
		},
	})
}

func HookQueueServer(lifecycle fx.Lifecycle, server *asynq.Server, mux *asynq.ServeMux) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Run(mux); err != nil {
					panic(err)
				}
				fmt.Println("MuxServer is running")
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.Shutdown()
			return nil
		},
	})
}

func HookStartWorker(
	lifecycle fx.Lifecycle,
	ctx context.Context,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	cfg *config.Config,
	mailer *mailer.Mailer,
	publisher rabbitmq.IPublisher,
	echo *echo.Echo,
	notificationScheduleRepository contracts.NotificationScheduleRepository,
	mux *asynq.ServeMux,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
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

			discountWorker := workers.NewDiscountWorker(mux, ctx, log, workers.HandleDiscountWorker)

			go func() {
				discountWorker.Start(&discountBase)
			}()
			return nil
		},
	})
}
