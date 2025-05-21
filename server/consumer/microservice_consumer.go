package consumer

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/streadway/amqp"
	"github.com/vmdt/notification-worker/shared"
)

func HandleConsumeMicroserviceMessage(queueName string, msg amqp.Delivery, dependencies *shared.MicroserviceBase) error {
	var microserviceMsg DiscountMessage
	if err := json.Unmarshal(msg.Body, &microserviceMsg); err != nil {
		dependencies.Log.Errorf("Error unmarshaling message: %v", err)
	}

	payload, err := json.Marshal(microserviceMsg)
	if err != nil {
		return err
	}

	if microserviceMsg.TaskQueueId != "" {
		dependencies.AsynqInspector.DeleteTask("discount:send", microserviceMsg.TaskQueueId)
	}

	task := asynq.NewTask("discount:send", payload)

	timezone := os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Ho_Chi_Minh"
		dependencies.Log.Warnf("TIMEZONE is not set, using default value: %s", timezone)
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		dependencies.Log.Errorf("Error loading location: %v", err)
	}

	layout := "2006-01-02 15:04:05"
	dependencies.Log.Infof("Schedule at: %s", microserviceMsg.ScheduleAt)
	scheduledTime, err := time.ParseInLocation(layout, microserviceMsg.ScheduleAt, loc)
	if err != nil {
		dependencies.Log.Errorf("Error parsing schedule_at: %v", err)
		return err
	}
	info, err := dependencies.AsynqClient.Enqueue(task, asynq.ProcessAt(scheduledTime))
	if err != nil {
		return err
	}

	dependencies.Log.Infof("Enqueued task: %s with ID: %s", info.Type, info.ID)

	dependencies.Publisher.PublishMessage(map[string]interface{}{
		"discount_id":   microserviceMsg.DiscountId,
		"task_queue_id": info.ID,
	}, "node_service", "node_key")

	return nil
}
