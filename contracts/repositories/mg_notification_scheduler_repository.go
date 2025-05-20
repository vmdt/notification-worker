package repositories

import (
	"context"
	"time"

	"github.com/vmdt/notification-worker/contracts"
	"github.com/vmdt/notification-worker/model"
	"github.com/vmdt/notification-worker/pkg/logger"
	"github.com/vmdt/notification-worker/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	statusPending = "pending"
)

type MongoNotificationScheduleRepository struct {
	log            logger.ILogger
	cfg            *mongodb.MongoDbOptions
	db             *mongo.Client
	collectionName string
}

func NewMongoNotificationScheduleRepository(log logger.ILogger, cfg *mongodb.MongoDbOptions, db *mongo.Client) contracts.NotificationScheduleRepository {
	return &MongoNotificationScheduleRepository{
		log:            log,
		cfg:            cfg,
		db:             db,
		collectionName: "notification_schedules",
	}
}

func (m *MongoNotificationScheduleRepository) CreateNotificationSchedule(schedule *model.NotificationSchedule) string {
	// Get collection
	collection := m.db.Database(m.cfg.Database).Collection(m.collectionName)

	// Set timestamps
	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	// Insert document
	result, err := collection.InsertOne(context.Background(), schedule)
	if err != nil {
		m.log.Errorf("Error creating notification schedule: %v", err)
		return ""
	}

	// Set the ID from the inserted document
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		schedule.ID = oid
		return oid.Hex()
	}

	return ""
}

func (m *MongoNotificationScheduleRepository) GetNotificationScheduleByScheduledAt(scheduledAt string) ([]*model.NotificationSchedule, error) {
	// Parse the input time (up to minutes)
	layout := "2006-01-02T15:04:05"
	parsedTime, err := time.ParseInLocation(layout, scheduledAt, time.UTC)
	if err != nil {
		m.log.Errorf("Error parsing scheduled time: %v", err)
		return nil, err
	}

	// Truncate to minute
	startOfMinute := parsedTime.Truncate(time.Minute)
	endOfMinute := startOfMinute.Add(time.Minute)

	// Build filter to match time range within the same minute
	filter := bson.M{
		"scheduled_at": bson.M{
			"$gte": startOfMinute,
			"$lt":  endOfMinute,
		},
		"status": statusPending,
	}

	// Query MongoDB
	collection := m.db.Database(m.cfg.Database).Collection(m.collectionName)
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		m.log.Errorf("Error finding notification schedules: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Decode result
	var schedules []*model.NotificationSchedule
	if err := cursor.All(context.Background(), &schedules); err != nil {
		m.log.Errorf("Error decoding notification schedules: %v", err)
		return nil, err
	}

	return schedules, nil
}
