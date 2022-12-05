package database

// @TODO: Implement context related code

import (
	"digitalpaper/backend/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Configuration parameters
var noFilterCriteria = bson.D{}

type Database struct {
	app      *core.Application
	Posts    *mongo.Collection
	Users    *mongo.Collection
	Sessions *mongo.Collection
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
