package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/pkg/logger"

	redis2 "github.com/vmdt/notification-worker/pkg/redis"
)

func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}

func NewServer(config *redis2.RedisOptions, logger logger.ILogger) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password: config.Password,
		},
		asynq.Config{Concurrency: 10},
	)
}
