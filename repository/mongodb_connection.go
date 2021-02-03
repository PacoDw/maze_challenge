package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBConn creates a new mongo db client to connect with the database.
func NewMongoDBConn(connString string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		log.Panicf("wrong connection with MongoDB: %s", err)
	}

	if err := client.Connect(context.Background()); err != nil {
		log.Panicf("the MongoDB connection can be reached: %s", err)
	}

	return client
}
