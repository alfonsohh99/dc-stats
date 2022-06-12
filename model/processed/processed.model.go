package processedModel

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

type User struct {
	Score       uint64        `bson:"score"`
	ChannelData []ChannelData `bson:"channel_data"`
}

type Guild struct {
	ID              primitive.ObjectID `bson:"_id"`
	GuildID         string             `bson:"guild_id"`
	TopUsers        []UserScore        `bson:"top_users"`
	TopMessageUsers []UserScore        `bson:"top_message_users"`
	UserData        map[string]User    `bson:"user_data"`
	UserMessageData map[string]User    `bson:"user_message_data"`
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
