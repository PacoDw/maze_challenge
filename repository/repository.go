package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// nolint
// mongoService represents a type for each service created, so all the servies like
// device in device.go must implement this type.
type mongoService struct {
	db *mongo.Client
}

// MongoDBService represenst the MongoDBService that contains services created.
// Note: If you has been created a new service it must be listed in this struct.
type MongoDBService struct {
	Quadrant QuadrantMongoDBService
	Spot     SpotMongoDBService
}

// New creates a new MongoDBService with all services in it
// Note: If you has been created a new service it must be listed in this struct.
func New(db *mongo.Client) *MongoDBService {
	return &MongoDBService{
		Quadrant: &QuadrantService{db: db},
		Spot:     &SpotService{db: db},
	}
}

// NewByConnString creates a new MongoDBService using a connection string.
func NewByConnString(connString string) *MongoDBService {
	return New(NewMongoDBConn(connString))
}
