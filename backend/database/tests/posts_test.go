package tests

import (
	"context"
	"digitalpaper/backend/core"
	"digitalpaper/backend/core/logger"
	"digitalpaper/backend/database"
	"fmt"
	mdb "github.com/ONSdigital/dp-mongodb-in-memory"
	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type DatabaseTestSuite struct {
	suite.Suite
	dummyContext context.Context
	dbServerMock *mdb.Server
	dbHandlers   *database.Database
}

// Helpers

var dummyPost = core.Post{
	Id:     "1234",
	Title:  "Dummy title",
	Body:   "Dummy body",
	Date:   "000000",
	Author: "Dummy author",
}

var dummyPosts = []core.Post{
	{
		Id:     "0000",
		Title:  "Dummy title zero",
		Body:   "Dummy body zero",
		Date:   "000000",
		Author: "Dummy author zero",
	},
	{
		Id:     "1111",
		Title:  "Dummy title one",
		Body:   "Dummy body one",
		Date:   "000001",
		Author: "Dummy author one",
	},
	{
		Id:     "2222",
		Title:  "Dummy title two",
		Body:   "Dummy body two",
		Date:   "000002",
		Author: "Dummy author two",
	},
}

func (dts *DatabaseTestSuite) fillDatabaseWithPosts() {
	for _, post := range dummyPosts {
		err := dts.dbHandlers.CreatePost(&dts.dummyContext, &post)

		if err != nil {
			dts.T().Error("error - could not fill database with posts")
		}
	}
}

// Configuration variables
var selectAll = bson.D{} // No filter criteria

// Setup

func (dts *DatabaseTestSuite) SetupSuite() {
	fmt.Println("Setting up database test suite...")

	dts.dummyContext = context.Background()

	// Setup in-memory MongoDB server
	dbServer, err := mdb.StartWithOptions(dts.dummyContext, "5.0.2", mdb.WithPort(27017))
	if err != nil {
		assert.Fail(dts.T(), "error while starting dummy MongoDB server")
	}

	// Setup database handlers/functions
	dummyApp := core.Application{
		Log:            logger.New(),
		SessionManager: &scs.SessionManager{},
	}

	dts.dbServerMock = dbServer
	dts.dbHandlers, err = database.NewDatabase(&dummyApp, dts.dbServerMock.URI())

	// @TODO: Create database connection test?
	assert.Nil(dts.T(), err, "error while creating new database instance")

	fmt.Println("Setting up database test suite... COMPLETE")
}

func (dts *DatabaseTestSuite) TearDownSuite() {
	fmt.Println("Tearing down database test suite...")

	dts.dbServerMock.Stop(dts.dummyContext)

	fmt.Println("Tearing down database test suite... COMPLETE")
}

func (dts *DatabaseTestSuite) TearDownTest() {
	deleteResult, _ := dts.dbHandlers.Posts.DeleteMany(dts.dummyContext, selectAll)
	fmt.Printf("Deleted %d posts", deleteResult)
}

// Tests

func (dts *DatabaseTestSuite) Test_GivenEmptyPostsDatabase_WhenFetchingAllPosts_ThenPostDataShouldBeEmpty() {
	posts, err := dts.dbHandlers.GetAllPosts(&dts.dummyContext)
	assert.Nil(dts.T(), err, "error raised while fetching posts in empty database, but shouldn't be")
	assert.Zero(dts.T(), posts, "no posts should be present in empty database")
}

func (dts *DatabaseTestSuite) Test_GivenInvalidPostData_WhenCreatingPost_ThenErrorIsReturned() {
	err := dts.dbHandlers.CreatePost(&dts.dummyContext, nil)
	assert.NotNil(dts.T(), err, "should raise an error while passing nil post data, but isn't")
}

func (dts *DatabaseTestSuite) Test_GivenValidPostData_WhenFetchingAllPosts_ThenCorrectPostDataIsReturned() {
	assert.Nil(dts.T(), dts.dbHandlers.CreatePost(&dts.dummyContext, &dummyPost), "error while creating post")

	posts, err := dts.dbHandlers.GetAllPosts(&dts.dummyContext)
	assert.Nil(dts.T(), err, "error while getting all posts")

	expectedId := "1234"
	actualId := posts[0].Id
	assert.Equal(dts.T(), expectedId, actualId, "post Id's should be the same, but aren't")

	expectedTitle := "Dummy title"
	actualTitle := posts[0].Title
	assert.Equal(dts.T(), expectedTitle, actualTitle, "post titles should be the same, but aren't")

	expectedBody := "Dummy body"
	actualBody := posts[0].Body
	assert.Equal(dts.T(), expectedBody, actualBody, "post bodies should be the same, but aren't")

	expectedDate := "000000"
	actualDate := posts[0].Date
	assert.Equal(dts.T(), expectedDate, actualDate, "post dates should be the same, but aren't")

	expectedAuthor := "Dummy author"
	actualAuthor := posts[0].Author
	assert.Equal(dts.T(), expectedAuthor, actualAuthor, "post authors should be the same, but aren't")
}

func (dts *DatabaseTestSuite) Test_GivenValidPostId_WhenGettingPostById_ThenCorrectPostIsReturned() {
	dts.fillDatabaseWithPosts()

	actualPost, err := dts.dbHandlers.GetPostById(&dts.dummyContext, "0000")
	assert.Nil(dts.T(), err, "error while fetching post by ID")
	assert.Equal(dts.T(), dummyPosts[0], actualPost)
}

func (dts *DatabaseTestSuite) Test_GivenNonExistingPostId_WhenGettingPostById_ThenEmptyPostAndErrorIsReturned() {
	dts.fillDatabaseWithPosts()

	actualPost, err := dts.dbHandlers.GetPostById(&dts.dummyContext, "123456789")
	assert.NotNil(dts.T(), err, "an error should be raised while fetching a post with non-existing ID, but isn't")
	assert.Equal(dts.T(), core.Post{}, actualPost)

	actualPost, err = dts.dbHandlers.GetPostById(&dts.dummyContext, "*-_?")
	assert.NotNil(dts.T(), err, "an error should be raised while fetching a post with non-existing ID, but isn't")
	assert.Equal(dts.T(), core.Post{}, actualPost)
}

func (dts *DatabaseTestSuite) Test_GivenFilledPostDatabase_WhenUpdatingPost_ThenCorrectPostIsUpdated() {
	dts.fillDatabaseWithPosts()

	updatedDummyPost := dummyPosts[1]
	updatedDummyPost.Title = "Updated dummy title"
	updatedDummyPost.Body = "Updated dummy body"
	updatedDummyPost.Date = "121212"
	updatedDummyPost.Author = "Updated dummy author"

	err := dts.dbHandlers.UpdatePost(&dts.dummyContext, &updatedDummyPost)
	assert.Nil(dts.T(), err, "error while fetching updating post")

	actualPost, err := dts.dbHandlers.GetPostById(&dts.dummyContext, updatedDummyPost.Id)
	assert.Nil(dts.T(), err, "error while fetching updating post")
	assert.Equal(dts.T(), updatedDummyPost, actualPost)
}

func (dts *DatabaseTestSuite) Test_GivenFilledPostDatabase_WhenDeletingAPost_ThenCorrectPostDeleted() {
	dts.fillDatabaseWithPosts()

	err := dts.dbHandlers.DeletePost(&dts.dummyContext, dummyPosts[1].Id)
	assert.Nil(dts.T(), err, "error while fetching updating post")

	actualPosts, err := dts.dbHandlers.GetAllPosts(&dts.dummyContext)
	assert.Nil(dts.T(), err, "error while fetching updating post")

	expectedPostsCount := 2
	actualPostsCount := len(actualPosts)
	assert.Equal(dts.T(), expectedPostsCount, actualPostsCount)
}

func TestRunDatabase(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}
