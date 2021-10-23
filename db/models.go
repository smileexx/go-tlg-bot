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

type Post struct {
	Id     string   `bson:"id"`
	Images []string `bson:"images"`
	Gifs   []string `bson:"gifs"`
	Shown  bool     `bson:"shown"`
	Tags   []string `bson:"tags"`
}
