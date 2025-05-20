package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout  = 60 * time.Second
	maxConnIdleTime = 3 * time.Second
	minPoolSize     = 20
	maxPoolSize     = 300
)

func NewMongoDB(cfg *MongoDbOptions) (*mongo.Client, error) {
	var uriAddress string
	if cfg.Uri == "" {
		uriAddress = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
		)
	} else {
		uriAddress = cfg.Uri
	}

	opt := options.Client().ApplyURI(uriAddress).
		SetConnectTimeout(connectTimeout).
		SetMaxConnIdleTime(maxConnIdleTime).
		SetMinPoolSize(minPoolSize).
		SetMaxPoolSize(maxPoolSize)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return client, nil
}
