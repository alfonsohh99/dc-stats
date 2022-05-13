package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DATA COLLECTION MODELS----
type Guild struct {
	ID      primitive.ObjectID `bson:"_id"`
	GuildID string             `bson:"guild_id"`
	Users   map[string]User    `bson:"users"`
}

type ChannelActivity struct {
	Score       uint64 `bson:"score"`
	ChannelName string `bson:"channel_name"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	UserName     string             `bson:"user_name"`
	UserActivity map[string]uint64  `bson:"user_activity"`
}

// PROCESSED COLLECTION MODELS
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
	ChannelData []ChannelData `bson:"channelData"`
}

type ProcessedGuild struct {
	ID       primitive.ObjectID       `bson:"_id"`
	GuildID  string                   `bson:"guild_id"`
	TopUsers []UserScore              `bson:"top_users"`
	UserData map[string]ProcessedUser `bson:"user_data"`
}
