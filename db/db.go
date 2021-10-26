package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "tlg-bot-db"

var ctx = context.TODO()
var client *mongo.Client
var db *mongo.Database

func CreateConnection() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database(dbName)
}

func GetRandomImages() []OldImage {
	imagesCollection := db.Collection("images")
	// opt := options.Find()
	// // opt.SetSort(bson.D{{"_id", -1}})
	// opt.SetLimit(100)
	matchStage := bson.D{{"$match", bson.D{{"shown", bson.D{{"$in", bson.A{nil, false}}}}}}}
	matchSample := bson.D{{"$sample", bson.D{{"size", 50}}}}
	cursor, err := imagesCollection.Aggregate(ctx, mongo.Pipeline{matchStage, matchSample})
	if err != nil {
		log.Fatal(err)
	}
	var images []OldImage
	if err = cursor.All(ctx, &images); err != nil {
		log.Fatal(err)
	}
	fmt.Println(images)
	return images
}

func GetRandomPosts(limit int) []Post {
	imagesCollection := db.Collection("posts")
	matchStage := bson.D{{"$match", bson.D{{"media", bson.D{{"$elemMatch", bson.D{{"shown", false}}}}}}}}
	// matchStage := bson.D{{"$match", bson.D{{"shown", bson.D{{"$in", bson.A{nil, false}}}}}}}
	matchSample := bson.D{{"$sample", bson.D{{"size", limit}}}}
	cursor, err := imagesCollection.Aggregate(ctx, mongo.Pipeline{matchStage, matchSample})
	if err != nil {
		log.Fatal(err)
	}
	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		log.Fatal(err)
	}
	return posts
}

func InsertPosts(posts []Post) {
	postsCollection := db.Collection("posts")
	var docs []interface{}

	for _, post := range posts {
		doc, _ := bson.Marshal(post)
		docs = append(docs, doc)
	}
	// docs, err := bson.(posts)
	opt := options.InsertMany()
	opt.SetOrdered(false)
	res, err := postsCollection.InsertMany(ctx, docs, opt)
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}

func SelectPostGifs() {
	postsCollection := db.Collection("posts")
	// bdoc := bson.D{{"media", bson.D{{"$elemMatch", bson.D{{"shown", true},{"type", "gif"} }}}}}
	var bdoc interface{}
	err := bson.UnmarshalExtJSON([]byte(`{"media": {"$elemMatch":{"shown":false, "type": "gif"}} }`), true, &bdoc)
	cursor, err := postsCollection.Find(ctx, bdoc)
	if err != nil {
		log.Fatal(err)
	}
	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		log.Fatal(err)
	}
	log.Println(posts)
}

func SelectPostsByTag(tag string) ([]Post, error) {
	var posts []Post
	postsCollection := db.Collection("posts")
	filter := bson.D{{"tags", bson.D{
		{"$regex", primitive.Regex{Pattern: tag, Options: "i"}},
	}}}
	// bson.RawValue{Type: bsontype.Regex, Value: []byte(fmt.Sprintf(`{"$regex":/%s/i}`, tag))}
	// sel := fmt.Sprintf(`{ "tags" : { $regex : /%s/i } }`, tag)
	// var bdoc interface{}
	// err := bson.UnmarshalExtJSON([]byte(sel), true, &bdoc)
	//bdoc := bson.D{{"tags", val}}
	cursor, err := postsCollection.Find(ctx, filter)
	if err != nil {
		return posts, err
	}
	if err = cursor.All(ctx, &posts); err != nil {
		return posts, err
	}
	if len(posts) < 1 {
		return posts, errors.New("No data with tag " + tag)
	}
	return posts, nil
}

func SelectPostsById(postId string) (*Post, error) {
	postsCollection := db.Collection("posts")
	cursor, err := postsCollection.Find(ctx, bson.D{{"id", postId}})
	if err != nil {
		log.Fatal(err)
	}
	var posts []Post
	if err = cursor.All(ctx, &posts); err != nil {
		log.Fatal(err)
		return nil, err
	}
	if len(posts) < 1 {
		return nil, errors.New(fmt.Sprintf(`Post %s not found =(`, postId))
	}
	return &posts[0], nil
}
