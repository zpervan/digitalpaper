package tests

import (
	"digitalpaper/backend/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PostsTestSuite struct {
	DatabaseTestSuite
}

// Mock data
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

// Helper functions
func (pts *PostsTestSuite) fillDatabaseWithPosts() {
	for _, post := range dummyPosts {
		err := pts.dbHandlers.CreatePost(pts.dummyContext, &post)

		if err != nil {
			pts.T().Error("error - could not fill database with posts")
		}
	}
}

func (pts *PostsTestSuite) TearDownTest() {
	_, _ = pts.dbHandlers.Posts.DeleteMany(pts.dummyContext, selectAll)
}

// Tests
func (pts *PostsTestSuite) Test_GivenEmptyPostsDatabase_WhenFetchingAllPosts_ThenPostDataShouldBeEmpty() {
	posts, err := pts.dbHandlers.GetAllPosts(pts.dummyContext)
	assert.Nil(pts.T(), err, "error raised while fetching posts in empty database, but shouldn't be")
	assert.Zero(pts.T(), posts, "no posts should be present in empty database")
}

func (pts *PostsTestSuite) Test_GivenInvalidPostData_WhenCreatingPost_ThenErrorIsReturned() {
	err := pts.dbHandlers.CreatePost(pts.dummyContext, nil)
	assert.NotNil(pts.T(), err, "should raise an error while passing nil post data, but isn't")
}

func (pts *PostsTestSuite) Test_GivenValidPostData_WhenFetchingAllPosts_ThenCorrectPostDataIsReturned() {
	assert.Nil(pts.T(), pts.dbHandlers.CreatePost(pts.dummyContext, &dummyPost), "error while creating post")

	posts, err := pts.dbHandlers.GetAllPosts(pts.dummyContext)
	assert.Nil(pts.T(), err, "error while getting all posts")

	expectedId := "1234"
	actualId := posts[0].Id
	assert.Equal(pts.T(), expectedId, actualId, "post Id's should be the same, but aren't")

	expectedTitle := "Dummy title"
	actualTitle := posts[0].Title
	assert.Equal(pts.T(), expectedTitle, actualTitle, "post titles should be the same, but aren't")

	expectedBody := "Dummy body"
	actualBody := posts[0].Body
	assert.Equal(pts.T(), expectedBody, actualBody, "post bodies should be the same, but aren't")

	expectedDate := "000000"
	actualDate := posts[0].Date
	assert.Equal(pts.T(), expectedDate, actualDate, "post dates should be the same, but aren't")

	expectedAuthor := "Dummy author"
	actualAuthor := posts[0].Author
	assert.Equal(pts.T(), expectedAuthor, actualAuthor, "post authors should be the same, but aren't")
}

func (pts *PostsTestSuite) Test_GivenValidPostId_WhenGettingPostById_ThenCorrectPostIsReturned() {
	pts.fillDatabaseWithPosts()

	actualPost, err := pts.dbHandlers.GetPostById(pts.dummyContext, "0000")
	assert.Nil(pts.T(), err, "error while fetching post by ID")
	assert.Equal(pts.T(), dummyPosts[0], actualPost)
}

func (pts *PostsTestSuite) Test_GivenNonExistingPostId_WhenGettingPostById_ThenEmptyPostAndErrorIsReturned() {
	pts.fillDatabaseWithPosts()

	actualPost, err := pts.dbHandlers.GetPostById(pts.dummyContext, "123456789")
	assert.NotNil(pts.T(), err, "an error should be raised while fetching a post with non-existing ID, but isn't")
	assert.Equal(pts.T(), core.Post{}, actualPost)

	actualPost, err = pts.dbHandlers.GetPostById(pts.dummyContext, "*-_?")
	assert.NotNil(pts.T(), err, "an error should be raised while fetching a post with non-existing ID, but isn't")
	assert.Equal(pts.T(), core.Post{}, actualPost)
}

func (pts *PostsTestSuite) Test_GivenFilledPostDatabase_WhenUpdatingPost_ThenCorrectPostIsUpdated() {
	pts.fillDatabaseWithPosts()

	updatedDummyPost := dummyPosts[1]
	updatedDummyPost.Title = "Updated dummy title"
	updatedDummyPost.Body = "Updated dummy body"
	updatedDummyPost.Date = "121212"
	updatedDummyPost.Author = "Updated dummy author"

	err := pts.dbHandlers.UpdatePost(pts.dummyContext, &updatedDummyPost)
	assert.Nil(pts.T(), err, "error while updating post, but shouldn't be")

	actualPost, err := pts.dbHandlers.GetPostById(pts.dummyContext, updatedDummyPost.Id)
	assert.Nil(pts.T(), err, "error while fetching post by Id, but shouldn't be")
	assert.Equal(pts.T(), updatedDummyPost, actualPost)
}

func (pts *PostsTestSuite) Test_GivenFilledPostDatabase_WhenDeletingAPost_ThenCorrectPostDeleted() {
	pts.fillDatabaseWithPosts()

	err := pts.dbHandlers.DeletePost(pts.dummyContext, dummyPosts[1].Id)
	assert.Nil(pts.T(), err, "error while deleting post, but shouldn't be")

	actualPosts, err := pts.dbHandlers.GetAllPosts(pts.dummyContext)
	assert.Nil(pts.T(), err, "error while fetching all posts, but shouldn't be")

	expectedPostsCount := 2
	actualPostsCount := len(actualPosts)
	assert.Equal(pts.T(), expectedPostsCount, actualPostsCount)
}

func TestRunDatabase(t *testing.T) {
	suite.Run(t, new(PostsTestSuite))
}
