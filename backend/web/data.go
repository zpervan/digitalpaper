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

var postMock = Post{Title: "This is a test title", Body: "Some beautifully written text with meaningful content.", Date: "01/01/1970", Author: "John Doe"}
