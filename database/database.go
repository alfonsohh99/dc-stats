package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"vc-stats/config"
)

var (
	DataCollection      *mongo.Collection
	ProcessedCollection *mongo.Collection
)

func Start(ctx context.Context) {

	clientOptions := options.Client().ApplyURI("mongodb://" + config.DatabaseEndpoint + ":" + config.DatabasePort + "/").SetAuth(options.Credential{
		Username: config.DatabaseUser,
		Password: config.DatabasePassword})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DataCollection = client.Database("discord").Collection("vc")
	ProcessedCollection = client.Database("discord").Collection("vc-processed")
}
