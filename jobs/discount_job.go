package jobs

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/server/consumer"
)

func AddDiscountJob(msg consumer.DiscountMessage, a *asynq.Client, log logger.ILogger) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	task := asynq.NewTask("discount:send", payload)

	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Ho_Chi_Minh"
		log.Warnf("TIMEZONE is not set, using default value: %s", timezone)
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Errorf("Error loading location: %v", err)
	}

	layout := "2006-01-02 15:04:05"
	scheduledTime, err := time.ParseInLocation(layout, msg.ScheduleAt, loc)
	if err != nil {
		log.Errorf("Error parsing schedule_at: %v", err)
		return err
	}
	info, err := a.Enqueue(task, asynq.ProcessAt(scheduledTime))
	if err != nil {
		return err
	}

	log.Infof("Enqueued task: %s with ID: %s", info.Type, info.ID)
	return nil
}
