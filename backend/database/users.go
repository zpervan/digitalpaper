package database

import (
	"context"
	"fmt"

	"digitalpaper/backend/core"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func (db Database) CreateUser(ctx context.Context, user core.User) error {
	if user.IsEmpty() {
		return fmt.Errorf("user data is empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	_, err = db.Users.InsertOne(ctx, user)

	if err != nil {
		return err
	}

	db.app.Log.Info("new user created")
	return nil
}

func (db Database) GetUsers(ctx context.Context, limit int) ([]core.User, error) {
	var filterOptions *options.FindOptions

	if limit != -1 {
		filterOptions = options.Find().SetLimit(int64(limit))
	}

	cursor, err := db.Users.Find(ctx, noFilterCriteria, filterOptions)

	if err != nil {
		return []core.User{}, err
	}

	var results []core.User

	for cursor.Next(ctx) {
		singleResult := core.User{}

		err = cursor.Decode(&singleResult)

		if err != nil {
			return nil, err
		}

		results = append(results, singleResult)
	}

	if len(results) == 0 {
		db.app.Log.Warn("couldn't find any users in database")
	}

	return results, nil
}

func (db Database) GetUserByUsername(ctx context.Context, username string) (_ core.User, err error) {
	defer func() {
		if err != nil {
			db.app.Log.Error(fmt.Sprintf("could not get user by username \"%s\". reason: %s", username, err))
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

func (db Database) GetUserByMail(ctx context.Context, mail string) (_ core.User, err error) {
	defer func() {
		if err != nil {
			db.app.Log.Error(fmt.Sprintf("could not get user by mail \"%s\". reason: %s", mail, err))
		}
	}()

	filter := bson.M{"mail": mail}
	queryResult := db.Users.FindOne(ctx, filter)

	if queryResult.Err() != nil {
		return core.User{}, queryResult.Err()
	}

	var user core.User
	err = queryResult.Decode(&user)

	if err != nil {
		return core.User{}, err
	}

	return user, nil
}

// @TODO: Should we update user by user ID?
func (db Database) UpdateUser(ctx context.Context, user core.User) error {
	userExists, err := db.UserExists(ctx, user)

	if err != nil {
		return err
	}

	if !userExists {
		return fmt.Errorf("cannot update non-existing user")
	}

	filter := bson.D{{"username", user.Username}}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	update := bson.D{{"$set", user}}
	result, err := db.Users.UpdateOne(ctx, filter, update)

	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		db.app.Log.Warn("update of user \"" + user.Username + "\" was unsuccessful")
	} else {
		db.app.Log.Info("modified user \"" + user.Username + "\"")
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
		db.app.Log.Warn("deleting user with username \"" + username + "\" was unsuccessful")
	} else {
		db.app.Log.Info("deleted user with username \"" + username + "\"")
	}

	return nil
}

func (db Database) UserExists(ctx context.Context, user core.User) (bool, error) {
	filterByUsername := bson.D{{"username", user.Username}}
	resultByUsername, err := db.Users.CountDocuments(ctx, filterByUsername, nil)
	if err != nil {
		return false, err
	}

	filterByMail := bson.D{{"mail", user.Mail}}
	resultByMail, err := db.Users.CountDocuments(ctx, filterByMail, nil)
	if err != nil {
		return false, err
	}

	return (resultByUsername != 0) || (resultByMail != 0), nil
}
