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

func GetImages() {
	imagesCollection := db.Collection("images")
	opt := options.Find()
	// opt.SetSort(bson.D{{"_id", -1}})
	opt.SetLimit(2)
	cursor, err := imagesCollection.Find(ctx, bson.M{}, opt)
	if err != nil {
		log.Fatal(err)
	}
	var images []Image
	if err = cursor.All(ctx, &images); err != nil {
		log.Fatal(err)
	}
	fmt.Println(images)
}
