package web

// @TODO: Implement context related code

import (
	"context"
	"digitalpaper/backend/core/logger"
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Posts *mongo.Collection
	Users *mongo.Collection
}

func Connect(dbUrl string) (Database, error) {
	clientOptions := options.Client().ApplyURI(dbUrl)
	clientOptions.SetServerSelectionTimeout(3 * time.Second)
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

func (db Database) CreatePost(post *Post) error {
	_, err := db.Posts.InsertOne(context.TODO(), post)

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Created new post with title \"%s\" and ID \"%s\"", post.Title, post.Id))
	return nil
}

func (db Database) GetAllPosts(context *context.Context) ([]Post, error) {
	filter := bson.M{}
	cursor, err := db.Posts.Find(*context, filter)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not retrieve all posts. Reason: %s", err))
	}

	var queryResults []Post

	for cursor.Next(*context) {
		singleResult := Post{}

		err := cursor.Decode(&singleResult)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not retrieve task. Reason: %s", err))
		}

		queryResults = append(queryResults, singleResult)
	}

	if len(queryResults) == 0 {
		logger.Warn("No tasks available")
	}

	return queryResults, nil
}
