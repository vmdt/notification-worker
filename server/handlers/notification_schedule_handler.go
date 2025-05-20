package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/vmdt/notification-worker/contracts"
	"github.com/vmdt/notification-worker/model"
	"github.com/vmdt/notification-worker/pkg/logger"
)

type NotificationScheduleHandler struct {
	log        logger.ILogger
	repository contracts.NotificationScheduleRepository
}

func NewNotificationScheduleHandler(log logger.ILogger, repository contracts.NotificationScheduleRepository) *NotificationScheduleHandler {
	return &NotificationScheduleHandler{
		log:        log,
		repository: repository,
	}
}

type CreateNotificationScheduleRequest struct {
	UserID      string                 `json:"user_id" validate:"required"`
	From        string                 `json:"from" validate:"required"`
	To          string                 `json:"to" validate:"required"`
	Subject     string                 `json:"subject" validate:"required"`
	Template    string                 `json:"template"`
	Content     string                 `json:"content" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	ScheduledAt string                 `json:"scheduled_at" validate:"required"`
	MaxRetries  int                    `json:"max_retries" validate:"required"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func (h *NotificationScheduleHandler) CreateNotificationSchedule(c echo.Context) error {
	var req CreateNotificationScheduleRequest
	if err := c.Bind(&req); err != nil {
		h.log.Errorf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Parse scheduled time
	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		h.log.Errorf("Error parsing scheduled time: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid scheduled_at format. Expected RFC3339 format",
		})
	}

	// Create notification schedule
	schedule := &model.NotificationSchedule{
		UserID:      req.UserID,
		From:        req.From,
		To:          req.To,
		Subject:     req.Subject,
		Template:    req.Template,
		Content:     req.Content,
		Type:        req.Type,
		Status:      "pending",
		ScheduledAt: scheduledAt,
		RetryCount:  0,
		MaxRetries:  req.MaxRetries,
		Metadata:    req.Metadata,
	}

	// Save to database
	id := h.repository.CreateNotificationSchedule(schedule)
	if id == "" {
		h.log.Error("Failed to create notification schedule")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create notification schedule",
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"id": id,
	})
}

func (h *NotificationScheduleHandler) GetNotificationByScheduleAt(c echo.Context) error {
	// Get schedule_at from query parameter
	scheduledAt := c.QueryParam("schedule_at")
	if scheduledAt == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "schedule_at query parameter is required",
		})
	}

	// Get notifications by scheduled time
	schedules, err := h.repository.GetNotificationScheduleByScheduledAt(scheduledAt)
	if err != nil {
		h.log.Errorf("Error getting notification schedules: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get notification schedules",
		})
	}

	return c.JSON(http.StatusOK, schedules)
}
