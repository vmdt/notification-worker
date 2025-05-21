package shared

import (
	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/rabbitmq"
)

type MicroserviceBase struct {
	Log            logger.ILogger
	Cfg            *config.Config
	AsynqClient    *asynq.Client
	AsynqInspector *asynq.Inspector
	Publisher      rabbitmq.IPublisher
}
