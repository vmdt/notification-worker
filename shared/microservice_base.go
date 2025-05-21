package shared

import (
	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/config"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type MicroserviceBase struct {
	Log         logger.ILogger
	Cfg         *config.Config
	AsynqClient *asynq.Client
}
