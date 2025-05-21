package workers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/jobs"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type DiscountWorker struct {
	queueName string
	mux       *asynq.ServeMux
	ctx       context.Context
	log       logger.ILogger
}

func NewDiscountWorker(mux *asynq.ServeMux, ctx context.Context, log logger.ILogger) *DiscountWorker {
	return &DiscountWorker{
		queueName: "discount:send",
		mux:       mux,
		ctx:       ctx,
		log:       log,
	}
}

func (d *DiscountWorker) Start() {
	d.log.Infof("Starting discount worker: %s", d.queueName)

	d.mux.HandleFunc(d.queueName, func(ctx context.Context, task *asynq.Task) error {
		var msg jobs.DiscountMessage
		if err := json.Unmarshal(task.Payload(), &msg); err != nil {
			d.log.Errorf("Error unmarshaling task payload: %v", err)
			return err
		}

		d.log.Infof("Processing discount message: %+v", msg)
		return nil
	})
}
