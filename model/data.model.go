package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Guild struct {
	ID      primitive.ObjectID `bson:"_id"`
	GuildID string             `bson:"guild_id"`
	Users   map[string]User    `bson:"users"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	UserName     string             `bson:"user_name"`
	UserActivity map[string]uint64  `bson:"user_activity"`
}
