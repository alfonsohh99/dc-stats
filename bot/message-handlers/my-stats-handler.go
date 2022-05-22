package messagehandlers

import (
	"context"
	"vc-stats/database"
	"vc-stats/model"
	"vc-stats/utils"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MyStats(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	optionsFindChannelData := options.FindOne().SetProjection(bson.M{
		"user_data." + m.Author.ID + ".score":        1,
		"user_data." + m.Author.ID + ".channel_data": bson.D{{Key: "$slice", Value: 10}},
	})
	findChannelData := database.ProcessedCollection.FindOne(ctx, filter, optionsFindChannelData)
	if findChannelData.Err() != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	findChannelData.Decode(&guildObject)

	if guildObject.UserData[m.Author.ID].ChannelData == nil || len(guildObject.UserData[m.Author.ID].ChannelData) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for you")
		return
	}

	var stats string
	stats += ":beginner: YOUR TOP CHANNELS :beginner:\n\n"

	for _, value := range guildObject.UserData[m.Author.ID].ChannelData {
		stats += value.ChannelName + ": " + utils.FormatTime(value.Score) + "\n"
	}
	stats += "\nTotal time: " + utils.FormatTime(guildObject.UserData[m.Author.ID].Score) + "\n"

	_, _ = s.ChannelMessageSend(m.ChannelID, stats)
}
