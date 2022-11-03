package web

// @TODO: Implement context related code

import (
	"context"
	"digitalpaper/backend/core/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Posts *mongo.Collection
	Users *mongo.Collection
}

func Connect(dbUrl string) (Database, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return Database{}, err
	}

	// Check if the database connection is alive
	if err := client.Ping(context.TODO(), nil); err != nil {
		return Database{}, err
	}

	// @TODO: Extract names in .env file
	postsCollection := client.Database("digital_paper").Collection("posts")
	usersCollection := client.Database("digital_paper").Collection("users")

	logger.Info("Database connection established")

	return Database{Posts: postsCollection, Users: usersCollection}, nil
}
