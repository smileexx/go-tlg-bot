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

type Schedule struct {
	Type string      `bson:"type"`
	Data interface{} `bson:"data"`
}
