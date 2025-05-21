package workers

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type DiscountWorker[T any] struct {
	queueName string
	mux       *asynq.ServeMux
	ctx       context.Context
	log       logger.ILogger
	handler   func(msg interface{}, dependencies T) error
}

func NewDiscountWorker[T any](
	mux *asynq.ServeMux,
	ctx context.Context,
	log logger.ILogger,
	handler func(msg interface{}, dependencies T) error,
) *DiscountWorker[T] {
	return &DiscountWorker[T]{
		queueName: "discount:send",
		mux:       mux,
		ctx:       ctx,
		log:       log,
		handler:   handler,
	}
}

func (d *DiscountWorker[T]) Start(dependencies T) {
	d.mux.HandleFunc(d.queueName, func(ctx context.Context, task *asynq.Task) error {
		d.handler(task.Payload(), dependencies)
		return nil
	})
}
