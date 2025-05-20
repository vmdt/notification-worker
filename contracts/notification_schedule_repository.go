package contracts

import (
	"github.com/vmdt/notification-worker/model"
)

type NotificationScheduleRepository interface {
	GetNotificationScheduleByScheduledAt(scheduledAt string) ([]*model.NotificationSchedule, error)
	CreateNotificationSchedule(*model.NotificationSchedule) string
}
