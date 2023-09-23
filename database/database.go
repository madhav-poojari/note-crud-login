package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var DbName = "note-making"
var UsersCollectionName = "users"
var NotesCollectionName = "notes"

func InitializeMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://arivappa:arivappa@hospital-management.qjwysaz.mongodb.net/")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	MongoClient = client
}
