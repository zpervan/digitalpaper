package core

type Post struct {
	Id     string `json:"id" bson:"id"`
	Title  string `json:"title" bson:"title"`
	Body   string `json:"body" bson:"body"`
	Date   string `json:"date" bson:"date"`
	Author string `json:"author" bson:"author"`
}

type User struct {
	Id       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Name     string `json:"name" bson:"name"`
	Surname  string `json:"surname" bson:"surname"`
	Mail     string `json:"mail" bson:"mail"`
	Password string `json:"password" bson:"password"`
}

type Comment struct {
	Id           string
	Author       string
	Date         string
	ParentPostId string
	Reply        string
}

func (u *User) IsEmpty() bool {
	isEmpty := true

	isEmpty = isEmpty && (u.Username == "")
	isEmpty = isEmpty && (u.Name == "")
	isEmpty = isEmpty && (u.Surname == "")
	isEmpty = isEmpty && (u.Mail == "")
	isEmpty = isEmpty && (u.Password == "")

	return isEmpty
}

var postMock = Post{Title: "This is a test title", Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean ultricies, tortor consectetur rhoncus sodales, arcu eros sollicitudin est, ac viverra lorem lectus at ligula. Nulla tempus tortor purus. In eget est bibendum, accumsan urna ultrices, porta eros. Cras ut congue lorem, non pulvinar dolor. Quisque sit amet sapien malesuada, aliquam purus quis, aliquam ligula. Etiam tristique, tortor vel tincidunt lacinia, est tortor pellentesque neque, ac pretium nisl sapien quis elit. Vestibulum eleifend velit fringilla pellentesque euismod. Duis fringilla dapibus velit, id mattis leo facilisis sit amet. Sed interdum euismod finibus.", Date: "01/01/1970", Author: "John Doe"}
var postMocks = []Post{postMock, {Title: "This is a new test title", Body: "Some ugly written text with cringe content.", Date: "01/01/1969", Author: "Booby Hobby"}}
