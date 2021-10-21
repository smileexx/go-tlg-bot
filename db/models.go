package db

type Content struct {
	Type   string `bson:"type"`
	FileId string `bson:"file_id"`
	Name   string `bson:"name"`
}
