package database

import (
	"context"
	"dc-stats/api/api-model"
	"dc-stats/config"
	"dc-stats/model/data"
	"dc-stats/model/processed"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DataCollection      *mongo.Collection
	ProcessedCollection *mongo.Collection
	ApiUsersCollection  *mongo.Collection
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
	ApiUsersCollection = client.Database("discord").Collection("api-users")
}

func FindOrCreateDataGuild(id string, ctx context.Context) (guild dataModel.Guild, err error) {

	var guildObject dataModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: id}}
	findGuild := DataCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		//Guild doesnt exist
		log.Println("GUILD NOT PRESENT", findGuild.Err().Error())
		newGuild := dataModel.CreateGuild(id)
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

func FindDataGuild(ctx context.Context, guildId string) (dataModel.Guild, error) {
	var guildObject dataModel.Guild
	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findGuild := DataCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		//Guild doesnt have data yet
		return guildObject, findGuild.Err()
	}

	findGuild.Decode(&guildObject)
	return guildObject, nil
}

func FindProcessedGuild(ctx context.Context, guildId string) (processedModel.Guild, error) {
	var guildObject processedModel.Guild
	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findGuild := ProcessedCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		//Guild doesnt have data yet
		return guildObject, findGuild.Err()
	}

	findGuild.Decode(&guildObject)
	return guildObject, nil
}

func UpdateDataGuildUsers(guildObject dataModel.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"users", guildObject.Users}}},
	})
}

func UpdateDataGuildUsersAndChannelMarks(guildObject dataModel.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"channel_marks", guildObject.ChannelMarks}}},
		{"$set", bson.D{{"users", guildObject.Users}}},
	})
}

func UpdateDataGuildUserNicknameMap(guildObject dataModel.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"user_nickname_map", guildObject.UserNicknameMap}}},
	})
}

func UpdateDataGuildChannelNameMap(guildObject dataModel.Guild, ctx context.Context) {
	DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
		{"$set", bson.D{{"channel_name_map", guildObject.ChannelNameMap}}},
	})
}

func SaveOrUpdateProcessedGuildFromVoice(guildId string, scores []processedModel.UserScore, userData map[string]processedModel.User, ctx context.Context) (processedModel.Guild, error) {

	var processedGuildObject processedModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findProcessedGuild := ProcessedCollection.FindOne(ctx, filter)
	if findProcessedGuild.Err() != nil {
		//Guild doesnt exist
		log.Println("PROCESSED GUILD NOT PRESENT", findProcessedGuild.Err().Error())
		newGuild := processedModel.CreateGuild(guildId)
		newGuild.TopUsers = scores
		newGuild.UserData = userData
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

func SaveOrUpdateProcessedGuildFromMessage(guildId string, scores []processedModel.UserScore, userData map[string]processedModel.User, ctx context.Context) (processedModel.Guild, error) {

	var processedGuildObject processedModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: guildId}}
	findProcessedGuild := ProcessedCollection.FindOne(ctx, filter)
	if findProcessedGuild.Err() != nil {
		//Guild doesnt exist
		log.Println("PROCESSED GUILD NOT PRESENT", findProcessedGuild.Err().Error())
		newGuild := processedModel.CreateGuild(guildId)
		newGuild.TopMessageUsers = scores
		newGuild.UserMessageData = userData
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
			{"$set", bson.D{{"top_message_users", scores}}},
			{"$set", bson.D{{"user_message_data", userData}}},
		})
	}
	return processedGuildObject, nil
}

func SaveOrUpdatApiUser(user apiModel.User, ctx context.Context) (apiModel.User, error) {

	var foundUserApi apiModel.User

	filter := bson.D{primitive.E{Key: "user_id", Value: user.UserId}}
	findApiUser := ApiUsersCollection.FindOne(ctx, filter)
	if findApiUser.Err() != nil {
		//User Api doesnt exist
		log.Println("API_USER NOT PRESENT: ", findApiUser.Err().Error())
		data, err := ApiUsersCollection.InsertOne(ctx, user)
		if err != nil {
			log.Println(err)
			return user, err
		}
		log.Println("CREATED API_USER:  ", data)
	} else {
		findApiUser.Decode(&foundUserApi)
		ApiUsersCollection.UpdateByID(ctx, foundUserApi.ID, user)
	}
	return user, nil
}

func FindApiUser(ctx context.Context, userId string) (apiModel.User, error) {
	var apiUser apiModel.User
	filter := bson.D{primitive.E{Key: "user_id", Value: userId}}
	findApiUser := ApiUsersCollection.FindOne(ctx, filter)
	if findApiUser.Err() != nil {
		//User not found
		return apiUser, findApiUser.Err()
	}
	findApiUser.Decode(&apiUser)
	return apiUser, nil
}
