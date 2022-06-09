package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserScore struct {
	Username string `bson:"user_name"`
	Score    uint64 `bson:"score"`
}

type ChannelData struct {
	Score       uint64 `bson:"score"`
	ChannelName string `bson:"channel_name"`
}

type ProcessedUser struct {
	Score       uint64        `bson:"score"`
	ChannelData []ChannelData `bson:"channel_data"`
}

type ProcessedGuild struct {
	ID              primitive.ObjectID       `bson:"_id"`
	GuildID         string                   `bson:"guild_id"`
	TopUsers        []UserScore              `bson:"top_users"`
	TopMessageUsers []UserScore              `bson:"top_message_users"`
	UserData        map[string]ProcessedUser `bson:"user_data"`
	UserMessageData map[string]ProcessedUser `bson:"user_message_data"`
}

func CreateProcessedGuild(id string) (guild ProcessedGuild) {
	return ProcessedGuild{
		ID:              primitive.NewObjectID(),
		GuildID:         id,
		TopUsers:        []UserScore{},
		TopMessageUsers: []UserScore{},
		UserData:        map[string]ProcessedUser{},
		UserMessageData: map[string]ProcessedUser{},
	}
}
