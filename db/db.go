package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "tlg-bot-db"

var ctx = context.TODO()
var client *mongo.Client
var db *mongo.Database

func Init() {
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

func GetImages() []Image {
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
	var images []Image
	if err = cursor.All(ctx, &images); err != nil {
		log.Fatal(err)
	}
	fmt.Println(images)
	return images
}
