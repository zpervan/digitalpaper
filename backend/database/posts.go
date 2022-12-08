package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"digitalpaper/backend/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	database := os.Getenv("MONGO_INITDB_DATABASE")

	db.Posts = client.Database(database).Collection("posts")
	db.Users = client.Database(database).Collection("users")
	db.Sessions = client.Database(database).Collection("sessions")

	db.app.Log.Info("Database connection established")

	return nil
}

func (db Database) CreatePost(ctx *context.Context, post *core.Post) error {
	_, err := db.Posts.InsertOne(context.TODO(), post)

	if err != nil {
		return err
	}

	db.app.Log.Info(fmt.Sprintf("Created new post with title \"%s\" and ID \"%s\"", post.Title, post.Id))
	return nil
}

func (db Database) GetAllPosts(context *context.Context) ([]core.Post, error) {
	filter := bson.M{}
	cursor, err := db.Posts.Find(*context, filter)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not retrieve all posts. Reason: %s", err))
	}

	var queryResults []core.Post

	for cursor.Next(*context) {
		singleResult := core.Post{}

		err := cursor.Decode(&singleResult)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not retrieve task. Reason: %s", err))
		}

		queryResults = append(queryResults, singleResult)
	}

	if len(queryResults) == 0 {
		db.app.Log.Warn("No tasks available")
	}

	return queryResults, nil
}

func (db Database) GetPostById(ctx *context.Context, id string) (_ core.Post, err error) {
	defer func() {
		if err != nil {
			db.app.Log.Error(fmt.Sprintf("could not fetch post by ID %s. Reason: %s", id, err))
		}
	}()

	filter := bson.M{"id": id}
	queryResult := db.Posts.FindOne(*ctx, filter)

	err = queryResult.Err()
	if err != nil {
		return core.Post{}, err
	}

	var result core.Post
	err = queryResult.Decode(&result)

	if err != nil {
		return core.Post{}, err
	}

	return result, nil
}

func (db Database) UpdatePost(ctx context.Context, updatedPost *core.Post) error {
	filter := bson.D{{"id", updatedPost.Id}}
	update := bson.D{{"$set", updatedPost}}

	result, err := db.Posts.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		db.app.Log.Warn("Update of post with Id " + updatedPost.Id + " was unsuccessful")
	} else {
		db.app.Log.Info("Modified post with Id " + updatedPost.Id)
	}

	return nil
}

func (db Database) DeletePost(ctx context.Context, postId string) error {
	filter := bson.D{{"id", postId}}

	result, err := db.Posts.DeleteOne(ctx, filter, nil)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		db.app.Log.Warn("Deleting post with Id " + postId + " was unsuccessful")
	} else {
		db.app.Log.Info("Deleted post with Id " + postId)
	}

	return nil
}