package web

type Post struct {
	Id     string `json:"id" bson:"id"`
	Title  string `json:"title" bson:"title"`
	Body   string `json:"body" bson:"body"`
	Date   string `json:"date" bson:"date"`
	Author string `json:"author" bson:"author"`
}

type Comment struct {
	Id           string
	Author       string
	Date         string
	ParentPostId string
	Reply        string
}

var postMock = Post{Title: "This is a test title", Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean ultricies, tortor consectetur rhoncus sodales, arcu eros sollicitudin est, ac viverra lorem lectus at ligula. Nulla tempus tortor purus. In eget est bibendum, accumsan urna ultrices, porta eros. Cras ut congue lorem, non pulvinar dolor. Quisque sit amet sapien malesuada, aliquam purus quis, aliquam ligula. Etiam tristique, tortor vel tincidunt lacinia, est tortor pellentesque neque, ac pretium nisl sapien quis elit. Vestibulum eleifend velit fringilla pellentesque euismod. Duis fringilla dapibus velit, id mattis leo facilisis sit amet. Sed interdum euismod finibus.", Date: "01/01/1970", Author: "John Doe"}
var postMocks = []Post{postMock, {Title: "This is a new test title", Body: "Some ugly written text with cringe content.", Date: "01/01/1969", Author: "Booby Hobby"}}
