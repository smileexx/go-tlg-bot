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

const (
	dbName     = "tlg-bot-db"
	clSchedule = "schedule"
	clBoobs    = "posts"
	clMemes    = "memes"
)

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

func GetRandomPosts(limit int) []Post {
	imagesCollection := db.Collection(clBoobs)
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
	postsCollection := db.Collection(clBoobs)
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
	postsCollection := db.Collection(clBoobs)
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
	postsCollection := db.Collection(clBoobs)
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
	postsCollection := db.Collection(clBoobs)
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

func UpdatePost(post Post) error {
	collection := db.Collection(clBoobs)
	filter := bson.D{{"id", post.Id}}
	data, err := bson.Marshal(post)
	if err != nil {
		return err
	}
	var doc *bson.D
	err = bson.Unmarshal(data, &doc)
	update := bson.D{{"$set", doc}}
	opt := options.UpdateOptions{}
	opt.SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, &opt)
	return err
}

func SaveSchedule(schedule Schedule) error {
	collection := db.Collection(clSchedule)
	filter := bson.D{{"chat_id", schedule.ChatId}, {"type", schedule.Type}}
	data, err := bson.Marshal(schedule)
	if err != nil {
		return err
	}
	var doc *bson.D
	err = bson.Unmarshal(data, &doc)
	update := bson.D{{"$set", doc}}
	opt := options.UpdateOptions{}
	opt.SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, &opt)
	return err
}

func RemoveFromSchedule(chatId int, feedType string) error {
	collection := db.Collection(clSchedule)
	filter := bson.D{{"chat_id", chatId}, {"type", feedType}}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

func SelectSchedule() ([]Schedule, error) {
	collection := db.Collection(clSchedule)
	var schedules []Schedule
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return schedules, err
	}
	if err = cursor.All(ctx, &schedules); err != nil {
		return schedules, err
	}
	if len(schedules) < 1 {
		return schedules, errors.New("No schedules")
	}
	return schedules, nil
}

func InsertMemesPosts(posts []MemePost) error {
	collection := db.Collection(clMemes)
	var docs []interface{}

	for _, post := range posts {
		doc, _ := bson.Marshal(post)
		docs = append(docs, doc)
	}
	opt := options.InsertMany()
	opt.SetOrdered(false)
	res, err := collection.InsertMany(ctx, docs, opt)
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

func GetRandomMemes(limit int) []MemePost {
	collection := db.Collection(clMemes)
	matchStage := bson.D{{"$match", bson.D{{"shown", false}}}}
	matchSample := bson.D{{"$sample", bson.D{{"size", limit}}}}
	cursor, err := collection.Aggregate(ctx, mongo.Pipeline{matchStage, matchSample})
	if err != nil {
		log.Fatal(err)
	}
	var posts []MemePost
	if err = cursor.All(ctx, &posts); err != nil {
		return posts
	}
	return posts
}

func UpdateMeme(meme MemePost) error {
	collection := db.Collection(clMemes)
	filter := bson.D{{"id", meme.Id}}
	data, err := bson.Marshal(meme)
	if err != nil {
		return err
	}
	var doc *bson.D
	err = bson.Unmarshal(data, &doc)
	update := bson.D{{"$set", doc}}
	opt := options.UpdateOptions{}
	opt.SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, &opt)
	return err
}
