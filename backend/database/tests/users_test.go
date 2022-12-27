package tests

import (
	"digitalpaper/backend/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UsersTestSuite struct {
	DatabaseTestSuite
}

// Configuration variables
const fetchAllUsers = -1

// Mock data
var dummyUser = core.User{
	Id:       "1234-abcd",
	Username: "dummyuser",
	Name:     "Dummyname",
	Surname:  "Dummysurname",
	Mail:     "dummy@mail.com",
	Password: "dummypassword",
}

var dummyUsers = []core.User{
	{
		Id:       "0000-abcd",
		Username: "dummyuser0",
		Name:     "Dummyname0",
		Surname:  "Dummysurname0",
		Mail:     "dummy0@mail.com",
		Password: "dummypassword0",
	},
	{
		Id:       "0010-abcd",
		Username: "dummyuser1",
		Name:     "Dummyname1",
		Surname:  "Dummysurname1",
		Mail:     "dummy1@mail.com",
		Password: "dummypassword1",
	},
	{
		Id:       "0200-abcd",
		Username: "dummyuser2",
		Name:     "Dummyname2",
		Surname:  "Dummysurname2",
		Mail:     "dummy2@mail.com",
		Password: "dummypassword2",
	},
}

// Helper functions
func (uts *UsersTestSuite) fillDatabaseWithUsers() {
	for _, user := range dummyUsers {
		err := uts.dbHandlers.CreateUser(uts.dummyContext, user)

		if err != nil {
			uts.T().Error("error - could not fill database with users")
		}
	}
}

func (uts *UsersTestSuite) checkUser(expectedUser core.User, actualUser core.User) bool {
	assert.Equal(uts.T(), expectedUser.Id, actualUser.Id, "user Id's should be the same, but aren't")
	assert.Equal(uts.T(), expectedUser.Username, actualUser.Username, "user usernames should be the same, but aren't")
	assert.Equal(uts.T(), expectedUser.Name, actualUser.Name, "user names should be the same, but aren't")
	assert.Equal(uts.T(), expectedUser.Surname, actualUser.Surname, "user surnames should be the same, but aren't")
	assert.Equal(uts.T(), expectedUser.Mail, actualUser.Mail, "user mails should be the same, but aren't")

	// Check only if password exists as it's hashed
	assert.NotEmpty(uts.T(), actualUser.Password)

	return true
}

func (uts *UsersTestSuite) TearDownTest() {
	_, _ = uts.dbHandlers.Users.DeleteMany(uts.dummyContext, selectAll)
}

// Tests
func (uts *UsersTestSuite) Test_GivenEmptyUsersDatabase_WhenFetchingUsers_ThenUserDataShouldBeEmpty() {
	users, err := uts.dbHandlers.GetUsers(uts.dummyContext, fetchAllUsers)
	assert.Nil(uts.T(), err, "error raised while fetching users in empty database, but shouldn't be")
	assert.Zero(uts.T(), users, "no users should be present in empty database")
}

func (uts *UsersTestSuite) Test_GivenEmptyUserData_WhenCreatingUser_ThenErrorIsReturned() {
	err := uts.dbHandlers.CreateUser(uts.dummyContext, core.User{})
	assert.NotNil(uts.T(), err, "should raise an error while passing nil user data, but isn't")
}

func (uts *UsersTestSuite) Test_GivenValidUserData_WhenFetchingUser_ThenCorrectUserDataIsReturned() {
	assert.Nil(uts.T(), uts.dbHandlers.CreateUser(uts.dummyContext, dummyUser), "error while creating user, but shouldn't be")

	users, err := uts.dbHandlers.GetUsers(uts.dummyContext, fetchAllUsers)
	assert.Nil(uts.T(), err, "error while getting user, but shouldn't be")

	actualId := users[0].Id
	assert.Equal(uts.T(), dummyUser.Id, actualId, "user Id's should be the same, but aren't")

	actualUsername := users[0].Username
	assert.Equal(uts.T(), dummyUser.Username, actualUsername, "user usernames should be the same, but aren't")

	actualName := users[0].Name
	assert.Equal(uts.T(), dummyUser.Name, actualName, "user names should be the same, but aren't")

	actualSurname := users[0].Surname
	assert.Equal(uts.T(), dummyUser.Surname, actualSurname, "user surnames should be the same, but aren't")

	actualMail := users[0].Mail
	assert.Equal(uts.T(), dummyUser.Mail, actualMail, "user mails should be the same, but aren't")

	// Password should be hashed so just compare that it's filled and not equal to the dummy user password
	assert.NotNil(uts.T(), users[0].Password)
	assert.NotEmpty(uts.T(), users[0].Password)

	actualPassword := users[0].Password
	assert.NotEqual(uts.T(), dummyUser.Password, actualPassword, "user passwords should be the same, but aren't")
}

func (uts *UsersTestSuite) Test_GivenValidUserMail_WhenGettingUserByMail_ThenCorrectUserIsReturned() {
	uts.fillDatabaseWithUsers()

	actualUser, err := uts.dbHandlers.GetUserByMail(uts.dummyContext, dummyUsers[1].Mail)
	assert.Nil(uts.T(), err, "error while fetching user by mail, but shouldn't be")
	uts.checkUser(dummyUsers[1], actualUser)
}

func (uts *UsersTestSuite) Test_GivenNonExistingUserMail_WhenGettingUserByMail_ThenEmptyUserAndErrorIsReturned() {
	uts.fillDatabaseWithUsers()

	actualUser, err := uts.dbHandlers.GetUserByMail(uts.dummyContext, "doesntexist@mail.com")
	assert.NotNil(uts.T(), err, "an error should be raised while fetching a user with non-existing mail, but isn't")
	assert.Equal(uts.T(), core.User{}, actualUser)
}

func (uts *UsersTestSuite) Test_GivenValidUserUsername_WhenGettingUserByUsername_ThenCorrectUserIsReturned() {
	uts.fillDatabaseWithUsers()

	actualUser, err := uts.dbHandlers.GetUserByUsername(uts.dummyContext, dummyUsers[1].Username)
	assert.Nil(uts.T(), err, "error while fetching user by username, but shouldn't be")
	uts.checkUser(dummyUsers[1], actualUser)
}

func (uts *UsersTestSuite) Test_GivenNonExistingUserUsername_WhenGettingUserByUsername_ThenEmptyUserAndErrorIsReturned() {
	uts.fillDatabaseWithUsers()

	actualUser, err := uts.dbHandlers.GetUserByUsername(uts.dummyContext, "doesntexist")
	assert.NotNil(uts.T(), err, "an error should be raised while fetching a user with non-existing username, but isn't")
	assert.Equal(uts.T(), core.User{}, actualUser)
}

func (uts *UsersTestSuite) Test_GivenFilledUserDatabase_WhenUpdatingUser_ThenCorrectUserIsUpdated() {
	uts.fillDatabaseWithUsers()

	updatedDummyUser := dummyUsers[1]
	updatedDummyUser.Name = "Updated dummy name"
	updatedDummyUser.Surname = "Updated dummy surname"
	updatedDummyUser.Mail = "newdummymail@mail.com"
	updatedDummyUser.Password = "newdummypassword"

	err := uts.dbHandlers.UpdateUser(uts.dummyContext, updatedDummyUser)
	assert.Nil(uts.T(), err, "error while updating user, but shouldn't be")

	actualUser, err := uts.dbHandlers.GetUserByMail(uts.dummyContext, updatedDummyUser.Mail)
	assert.Nil(uts.T(), err, "error while fetching user, but shouldn't be")

	uts.checkUser(updatedDummyUser, actualUser)
}

func (uts *UsersTestSuite) Test_GivenFilledUserDatabase_WhenUpdatingNonExistingUser_ThenErrorIsReturned() {
	uts.fillDatabaseWithUsers()

	nonExistingUser := core.User{Username: "doesntexits"}

	assert.Error(uts.T(), uts.dbHandlers.UpdateUser(uts.dummyContext, nonExistingUser), "an error should be raised as the user doesn't exist, but isn't")
}

func (uts *UsersTestSuite) Test_GivenFilledUserDatabase_WhenDeletingUser_ThenCorrectUserDeleted() {
	uts.fillDatabaseWithUsers()

	err := uts.dbHandlers.DeleteUser(uts.dummyContext, dummyUsers[1].Username)
	assert.Nil(uts.T(), err, "error while deleting user, but shouldn't")

	actualUsers, err := uts.dbHandlers.GetUsers(uts.dummyContext, fetchAllUsers)
	assert.Nil(uts.T(), err, "error while fetching users, but shouldn't")

	expectedUsersCount := 2
	actualUsersCount := len(actualUsers)
	assert.Equal(uts.T(), expectedUsersCount, actualUsersCount)
}

func (uts *UsersTestSuite) Test_GivenValidUserData_WhenCheckingUserExistence_ThenUserExists() {
	uts.fillDatabaseWithUsers()

	actualExistenceResult, err := uts.dbHandlers.UserExists(uts.dummyContext, dummyUsers[1])
	assert.Nil(uts.T(), err, "error while checking user existence, but shouldn't")
	assert.True(uts.T(), actualExistenceResult)
}

func (uts *UsersTestSuite) Test_GivenValidNonExistingUserData_WhenCheckingUserExistence_ThenUserDoesNotExist() {
	uts.fillDatabaseWithUsers()

	// Only needed fields are populated
	nonExistingUser := core.User{
		Username: "iamnothing",
		Mail:     "idonotexist@mail.com",
	}

	actualExistenceResult, err := uts.dbHandlers.UserExists(uts.dummyContext, nonExistingUser)
	assert.Nil(uts.T(), err, "error while checking user existence, but shouldn't")
	assert.False(uts.T(), actualExistenceResult)
}

func TestUsersRun(t *testing.T) {
	suite.Run(t, new(UsersTestSuite))
}
