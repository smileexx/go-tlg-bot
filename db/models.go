package db

type Content struct {
	Type   string `bson:"type"`
	FileId string `bson:"file_id"`
	Name   string `bson:"name"`
}

type Image struct {
	Id     string `bson:"id"`
	SrcUrl string `bson:"src"`
	Shown  bool   `bson:"shown"`
	Tags   string `bson:"tags"`
}
