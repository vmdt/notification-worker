package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/contracts"
	echoserver "github.com/vmdt/notification-worker/pkg/echo"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
	"github.com/vmdt/notification-worker/server/endpoints"
	"github.com/vmdt/notification-worker/server/handlers"
	"go.uber.org/fx"
)

func RunServers(
	lc fx.Lifecycle,
	log logger.ILogger,
	e *echo.Echo,
	ctx context.Context,
	cfg *config.Config,
	repository contracts.NotificationScheduleRepository,
	publisher rabbitmq.IPublisher,
) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoserver.RunHttpServer(ctx, e, log, cfg.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			// Create notification schedule handler
			notificationHandler := handlers.NewNotificationScheduleHandler(log, repository, publisher)

			// Register routes
			endpoints.RegisterNotificationScheduleRoutes(e, notificationHandler)

			// Health check endpoint
			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Info("shutting down server")
			return nil
		},
	})

	return nil
}
