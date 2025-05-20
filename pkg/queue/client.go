package queue

import (
	"fmt"

	"github.com/hibiken/asynq"
	redis2 "github.com/vmdt/notification-worker/pkg/redis"
)

func NewClient(config *redis2.RedisOptions) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", config.Host, config.Port)})
}
