package database

// @TODO: Implement context related code

import (
	"context"
	"digitalpaper/backend/core"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration parameters
var noFilterCriteria = bson.D{}
const databaseName = "digital_paper"

type Database struct {
	app      *core.Application
	Posts    *mongo.Collection
	Users    *mongo.Collection
	Sessions *mongo.Collection
}

func (db *Database) Connect(dbUrl string) error {
	clientOptions := options.Client().ApplyURI(dbUrl)
	clientOptions.SetServerSelectionTimeout(3 * time.Second)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	// Check if the database connection is alive
	if err := client.Ping(context.TODO(), nil); err != nil {
		return err
	}

	// In a local dev environment (non-Docker), fetching the given environment variable 
	// could return in an empty string and would not properly initialize the database
	database := os.Getenv("MONGO_INITDB_DATABASE")
	if database == "" {
		database = databaseName
	}

	db.Posts = client.Database(database).Collection("posts")
	db.Users = client.Database(database).Collection("users")
	db.Sessions = client.Database(database).Collection("sessions")

	db.app.Log.Info("Database connection established")

	return nil
}

func NewDatabase(app *core.Application, dbUrl string) (*Database, error) {
	database := &Database{}
	database.app = app

	err := database.Connect(dbUrl)
	if err != nil {
		return nil, err
	}

	return database, nil
}

