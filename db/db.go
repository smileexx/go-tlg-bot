package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var ctx = context.TODO()

func Connect() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	contentCollection := client.Database("tlg-bot-db").Collection("content")
	cursor, err := contentCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var cont []Content
	if err = cursor.All(ctx, &cont); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cont)
}
