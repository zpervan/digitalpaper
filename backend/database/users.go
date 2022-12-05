package database

import (
	"context"
	"fmt"

	"digitalpaper/backend/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (db Database) CreateUser(ctx context.Context, user *core.User) error {
	_, err := db.Users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	db.app.Log.Info("New user created")
	return nil
}

func (db Database) GetUsers(ctx *context.Context, limit int) ([]core.User, error) {
	var filterOptions *options.FindOptions

	if limit != -1 {
		filterOptions = options.Find().SetLimit(int64(limit))
	}

	cursor, err := db.Users.Find(*ctx, noFilterCriteria, filterOptions)

	if err != nil {
		return []core.User{}, err
	}

	var results []core.User

	for cursor.Next(*ctx) {
		singleResult := core.User{}

		err = cursor.Decode(&singleResult)

		if err != nil {
			return nil, err
		}

		results = append(results, singleResult)
	}

	if len(results) == 0 {
		db.app.Log.Warn("Couldn't find any users in database")
	}

	return results, nil
}

func (db Database) GetUserByUsername(ctx context.Context, username string) (_ core.User, err error) {
	defer func() {
		if err != nil {
			db.app.Log.Error(fmt.Sprintf("could not get user by username \"%s\". Reason: %s", username, err))
		}
	}()

	filter := bson.M{"username": username}
	queryResult := db.Users.FindOne(ctx, filter)

	err = queryResult.Err()
	if err != nil {
		return core.User{}, err
	}

	var user core.User
	err = queryResult.Decode(&user)

	if err != nil {
		return core.User{}, err
	}

	return user, nil
}

func (db Database) UpdateUser(ctx context.Context, user *core.User) error {
	filter := bson.D{{"username", user.Username}}
	update := bson.D{{"$set", user}}

	result, err := db.Users.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		db.app.Log.Warn("Update of user \"" + user.Username + "\" was unsuccessful")
	} else {
		db.app.Log.Info("Modified user \"" + user.Username + "\"")
	}

	return nil
}

func (db Database) DeleteUser(ctx context.Context, username string) error {
	filter := bson.D{{"username", username}}

	result, err := db.Users.DeleteOne(ctx, filter, nil)

	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		db.app.Log.Warn("Deleting user with username \"" + username + "\" was unsuccessful")
	} else {
		db.app.Log.Info("Deleted user with username \"" + username + "\"")
	}

	return nil
}
