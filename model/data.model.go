package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Guild struct {
	ID              primitive.ObjectID     `bson:"_id"`
	GuildID         string                 `bson:"guild_id"`
	Users           map[string]User        `bson:"users"`
	ChannelMarks    map[string]MessageMark `bson:"channel_marks"`
	ChannelNameMap  map[string]string      `bson:"channel_name_map"`
	UserNicknameMap map[string]string      `bson:"user_nickname_map"`
}

type MessageMark struct {
	BeforeId      string `bson:"before_id"`
	AfterId       string `bson:"after_id"`
	TotalMessages uint64 `bson:"total_messages"`
}

type User struct {
	ID                  primitive.ObjectID `bson:"_id"`
	UserID              string             `bson:"user_id"`
	UserVoiceActivity   map[string]uint64  `bson:"user_activity"`
	UserMessageActivity map[string]uint64  `bson:"user_message_activity"`
}

func CreateDataGuild(id string) (guild Guild) {
	return Guild{
		ID:              primitive.NewObjectID(),
		GuildID:         id,
		Users:           map[string]User{},
		ChannelMarks:    map[string]MessageMark{},
		ChannelNameMap:  map[string]string{},
		UserNicknameMap: map[string]string{},
	}
}

func CreateDataUser(id string) (user User) {
	return User{
		ID:                  primitive.NewObjectID(),
		UserID:              id,
		UserVoiceActivity:   map[string]uint64{},
		UserMessageActivity: map[string]uint64{},
	}
}
