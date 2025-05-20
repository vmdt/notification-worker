package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationSchedule represents a scheduled notification in the system
type NotificationSchedule struct {
	ID          primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	UserID      string                 `bson:"user_id" json:"user_id" validate:"required"`
	From        string                 `bson:"from" json:"from" validate:"required"` // sender information
	To          string                 `bson:"to" json:"to" validate:"required"`     // recipient information
	Subject     string                 `bson:"subject" json:"subject" validate:"required"`
	Template    string                 `bson:"template" json:"template"`
	Content     string                 `bson:"content" json:"content" validate:"required"`
	Type        string                 `bson:"type" json:"type" validate:"required"`     // email, sms, push, etc.
	Status      string                 `bson:"status" json:"status" validate:"required"` // pending, sent, failed
	ScheduledAt time.Time              `bson:"scheduled_at" json:"scheduled_at" validate:"required"`
	SentAt      *time.Time             `bson:"sent_at,omitempty" json:"sent_at,omitempty"`
	RetryCount  int                    `bson:"retry_count" json:"retry_count"`
	MaxRetries  int                    `bson:"max_retries" json:"max_retries" validate:"required"`
	Metadata    map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt   time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time              `bson:"updated_at" json:"updated_at"`
}

// NewNotificationSchedule creates a new notification schedule
func NewNotificationSchedule(userID, title, content, notificationType string, scheduledAt time.Time) *NotificationSchedule {
	now := time.Now()
	return &NotificationSchedule{
		UserID:      userID,
		From:        "",
		To:          "",
		Subject:     title,
		Template:    "",
		Content:     content,
		Type:        notificationType,
		Status:      "pending",
		ScheduledAt: scheduledAt,
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// MarkAsSent marks the notification as sent
func (n *NotificationSchedule) MarkAsSent() {
	now := time.Now()
	n.Status = "sent"
	n.SentAt = &now
	n.UpdatedAt = now
}

// MarkAsFailed marks the notification as failed and increments retry count
func (n *NotificationSchedule) MarkAsFailed() {
	n.Status = "failed"
	n.RetryCount++
	n.UpdatedAt = time.Now()
}

// CanRetry checks if the notification can be retried
func (n *NotificationSchedule) CanRetry() bool {
	return n.Status == "failed" && n.RetryCount < n.MaxRetries
}

// ResetForRetry resets the notification for retry
func (n *NotificationSchedule) ResetForRetry() {
	n.Status = "pending"
	n.UpdatedAt = time.Now()
}
