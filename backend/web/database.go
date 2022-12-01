package web

// @TODO: Implement context related code

import (
	"context"
	"digitalpaper/backend/core/logger"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Configuration parameters
var noFilterCriteria = bson.D{}

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

func (db Database) createPost(ctx *context.Context, post *Post) error {
	_, err := db.Posts.InsertOne(context.TODO(), post)

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Created new post with title \"%s\" and ID \"%s\"", post.Title, post.Id))
	return nil
}

func (db Database) getAllPosts(context *context.Context) ([]Post, error) {
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

func (db Database) getPostById(ctx *context.Context, id string) (_ Post, err error) {
	defer func() {
		if err != nil {
			logger.Error(fmt.Sprintf("could not fetch post by ID. Reason: %s", err))
		}
	}()

	filter := bson.M{"id": id}
	queryResult := db.Posts.FindOne(*ctx, filter)

	err = queryResult.Err()
	if err != nil {
		return Post{}, err
	}

	var result Post
	err = queryResult.Decode(&result)

	if err != nil {
		return Post{}, err
	}

	return result, nil
}

func (db Database) updatePost(ctx context.Context, updatedPost *Post) error {
	filter := bson.D{{"id", updatedPost.Id}}
	update := bson.D{{"$set", updatedPost}}

	result, err := db.Posts.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		logger.Warn("Update of post with Id " + updatedPost.Id + " was unsuccessful")
	} else {
		logger.Info("Modified post with Id " + updatedPost.Id)
	}

	return nil
}

func (db Database) deletePost(ctx context.Context, postId string) error{
	filter := bson.D{{"id", postId}}

	result, err := db.Posts.DeleteOne(ctx, filter, nil)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		logger.Warn("Deleting post with Id " + postId + " was unsuccessful")
	} else {
        logger.Info("Deleted post with Id " + postId)
    }

	return nil
}

func (db Database) createUser(ctx context.Context, user *User) error {
	_, err := db.Users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	logger.Info("New user created")
	return nil
}

func (db Database) getUsers(ctx *context.Context, limit int) ([]User, error) {
    var filterOptions *options.FindOptions

    if limit != -1 {
        filterOptions = options.Find().SetLimit(int64(limit))
    }

    cursor, err := db.Users.Find(*ctx, noFilterCriteria, filterOptions)

    if err != nil {
        return []User{}, err
    }

    var results []User

    for cursor.Next(*ctx) {
        singleResult := User{}

        err = cursor.Decode(&singleResult)

        if err != nil {
            return nil, err
        }

        results = append(results, singleResult)
    }

    if len(results) == 0 {
        logger.Warn("Couldn't find any users in database")
    }

    return results, nil
}

func (db Database) updateUser(ctx context.Context, user *User) error {
	filter := bson.D{{"username", user.Username}}
	update := bson.D{{"$set", user}}

	result, err := db.Users.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		logger.Warn("Update of user \"" + user.Username + "\" was unsuccessful")
	} else {
		logger.Info("Modified user \"" + user.Username + "\"")
	}
	
	return nil
}

func (db Database) getUserByUsername(ctx context.Context, username string) (User, error) {
	filter := bson.M{"username": username}
	queryResult := db.Users.FindOne(ctx, filter)

	if queryResult.Err() != nil {
		logger.Warn("Didn't find any entry with the username \"" + username + "\"")
		return User{}, nil
	}

	var user User
	err := queryResult.Decode(&user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db Database) deleteUser(ctx context.Context, username string) error {
	filter := bson.D{{"username", username}}

	result, err := db.Users.DeleteOne(ctx, filter, nil)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		logger.Warn("Deleting user with username \"" + username + "\" was unsuccessful")
	} else {
        logger.Info("Deleted user with username \"" + username + "\"")
    }

	return nil
}
