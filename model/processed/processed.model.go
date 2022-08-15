package processedModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserScore struct {
	Username string `json:"user_name" bson:"user_name"`
	Score    uint64 `json:"score" bson:"score"`
}

type ChannelData struct {
	Score       uint64 `json:"score" bson:"score"`
	ChannelName string `json:"channel_name" bson:"channel_name"`
}

type User struct {
	Score       uint64        `json:"score" bson:"score"`
	ChannelData []ChannelData `json:"channel_data" bson:"channel_data"`
}

type Guild struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id" `
	GuildID         string             `json:"guild_id" bson:"guild_id"`
	TopUsers        []UserScore        `json:"top_users" bson:"top_users"`
	TopMessageUsers []UserScore        `json:"top_message_users" bson:"top_message_users"`
	UserData        map[string]User    `json:"user_data" bson:"user_data"`
	UserMessageData map[string]User    `json:"user_message_data" bson:"user_message_data"`
}

func CreateGuild(id string) (guild Guild) {
	return Guild{
		ID:              primitive.NewObjectID(),
		GuildID:         id,
		TopUsers:        []UserScore{},
		TopMessageUsers: []UserScore{},
		UserData:        map[string]User{},
		UserMessageData: map[string]User{},
	}
}
