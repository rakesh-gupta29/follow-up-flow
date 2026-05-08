// internal/database/database.go
package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
}

func New(uri string) (*DB, error) {
	// 1. Configure Client Options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(25).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	// 3. Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	return &DB{Client: client}, nil
}

// GetCollection is a helper to access a specific collection
func (db *DB) GetCollection(dbName, colName string) *mongo.Collection {
	return db.Client.Database(dbName).Collection(colName)
}

func (db *DB) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Client.Disconnect(ctx); err != nil {
		fmt.Printf("Error closing mongodb connection: %v\n", err)
	}
}
