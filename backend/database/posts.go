package database

import (
	"context"
	"fmt"

	"digitalpaper/backend/core"

	"go.mongodb.org/mongo-driver/bson"
)

func (db Database) CreatePost(ctx context.Context, post *core.Post) error {
	_, err := db.Posts.InsertOne(context.TODO(), post)

	if err != nil {
		return err
	}

	db.app.Log.Info(fmt.Sprintf("created new post with title \"%s\" and ID \"%s\"", post.Title, post.Id))
	return nil
}

func (db Database) GetAllPosts(ctx context.Context) ([]core.Post, error) {
	filter := bson.M{}
	cursor, err := db.Posts.Find(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("could not retrieve all posts. reason: %s", err)
	}

	var queryResults []core.Post

	for cursor.Next(ctx) {
		singleResult := core.Post{}

		err := cursor.Decode(&singleResult)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve task. reason: %s", err)
		}

		queryResults = append(queryResults, singleResult)
	}

	if len(queryResults) == 0 {
		db.app.Log.Warn("no tasks available")
	}

	return queryResults, nil
}

// @TODO: Add check that post ID should only contain letters, numbers and the "-" symbol (i.e. 1234-abcd-a1b2)?
func (db Database) GetPostById(ctx context.Context, id string) (_ core.Post, err error) {
	defer func() {
		if err != nil {
			db.app.Log.Error(fmt.Sprintf("could not fetch post by ID %s. reason: %s", id, err))
		}
	}()

	filter := bson.M{"id": id}
	queryResult := db.Posts.FindOne(ctx, filter)

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

// @TODO: Handle situation when the modified count is 0 with an error?
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

// @TODO: Handle situation when the deleted count is 0 with an error?
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
