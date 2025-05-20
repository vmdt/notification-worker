package endpoints

import (
	"github.com/labstack/echo/v4"
	"github.com/vmdt/notification-worker/server/handlers"
)

func RegisterNotificationScheduleRoutes(e *echo.Echo, handler *handlers.NotificationScheduleHandler) {
	// Create a group for notification schedule endpoints
	notificationGroup := e.Group("/api/v1/notifications")

	// Register routes
	notificationGroup.POST("/schedule", handler.CreateNotificationSchedule)
	notificationGroup.GET("/schedule", handler.GetNotificationByScheduleAt)
}
