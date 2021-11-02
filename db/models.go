package db

type Content struct {
	Type   string `bson:"type"`
	FileId string `bson:"file_id"`
	Name   string `bson:"name"`
}

type OldImage struct {
	Id     string `bson:"id"`
	SrcUrl string `bson:"src"`
	Shown  bool   `bson:"shown"`
	Tags   string `bson:"tags,omitempty"`
	PostId string `bson:"post,omitempty"`
}

type Media struct {
	Type  string `bson:"type"`
	Src   string `bson:"src"`
	Shown bool   `bson:"shown"`
}

type Post struct {
	Id    string   `bson:"id"`
	Media []Media  `bson:"media"`
	Tags  []string `bson:"tags"`
}

type MemePost struct {
	Id        string   `bson:"id"`
	Title     string   `bson:"title"`
	Type      string   `bson:"type"`
	Src       string   `bson:"src"`
	Permalink string   `bson:"permalink"`
	Shown     bool     `bson:"shown"`
	Created   int      `bson:"created"`
	Tags      []string `bson:"tags"`
}

type Schedule struct {
	Type       string      `bson:"type"`
	Period     int         `bson:"period"`
	ChatId     int         `bson:"chat_id"`
	Title      string      `bson:"title,omitempty"`
	LastAction int64       `bson:"last_action"`
	Data       interface{} `bson:"data,omitempty"`
}

type User struct {
	Id        int    `bson:"id"`
	FirstName string `bson:"first_name"`
	UserName  string `bson:"username"`
	IsAdmin   bool   `bson:"is_admin"`
}
