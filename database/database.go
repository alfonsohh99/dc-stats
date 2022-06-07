package database

import (
	"context"
	"dc-stats/config"
	"dc-stats/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func FindOrCreateDataGuild(id string, ctx context.Context) (guild model.Guild, err error) {

	var guildObject model.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: id}}
	findGuild := DataCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		//Guild doesnt exist
		log.Println("GUILD NOT PRESSENT", findGuild.Err().Error())
		newGuild := model.CreateDataGuild(id)
		data, err := DataCollection.InsertOne(ctx, newGuild)
		if err != nil {
			return guildObject, err
		}
		log.Println("CREATED GUILD, ", data)
		guildObject = newGuild
	} else {
		findGuild.Decode(&guildObject)
	}

	return guildObject, nil
}

func FindDataGuild(ctx context.Context, guildId string) (model.Guild, error) {
	var guildObject model.Guild
	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findGuild := DataCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		//Guild doesnt have data yet
		return guildObject, findGuild.Err()
	}

	findGuild.Decode(&guildObject)
	return guildObject, nil
}

func UpdateDataGuildUsers(guildObject model.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"users", guildObject.Users}}},
	})
}

func UpdateDataGuildUsersAndChannelMarks(guildObject model.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"channel_marks", guildObject.ChannelMarks}}},
		{"$set", bson.D{{"users", guildObject.Users}}},
	})
}

func SaveOrUpdateProcessedGuild(guildId string, scores []model.UserScore, userData map[string]model.ProcessedUser, ctx context.Context) (model.ProcessedGuild, error) {

	var processedGuildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findProcessedGuild := ProcessedCollection.FindOne(ctx, filter)
	if findProcessedGuild.Err() != nil {
		//Guild doesnt exist
		log.Println("PROCESSED GUILD NOT PRESSENT", findProcessedGuild.Err().Error())
		newGuild := model.ProcessedGuild{
			ID:       primitive.NewObjectID(),
			GuildID:  guildId,
			TopUsers: scores,
			UserData: userData,
		}
		data, err := ProcessedCollection.InsertOne(ctx, newGuild)
		if err != nil {
			log.Println(err)
			return processedGuildObject, err
		}
		log.Println("CREATED PROCESSED GUILD, ", data)
		processedGuildObject = newGuild
	} else {
		findProcessedGuild.Decode(&processedGuildObject)
		ProcessedCollection.UpdateByID(ctx, processedGuildObject.ID, bson.D{
			{"$set", bson.D{{"top_users", scores}}},
			{"$set", bson.D{{"user_data", userData}}},
		})
	}
	return processedGuildObject, nil
}
