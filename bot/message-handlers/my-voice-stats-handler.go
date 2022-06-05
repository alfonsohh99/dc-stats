package messagehandlers

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"dc-stats/utils"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MY_VOICE_STATS_MESSAGE_HEADING = ":beginner: YOUR TOP ACTIVE VOICE CHANNELS :beginner:\n\n"

func MyVoiceStats(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	optionsFindChannelData := options.FindOne().SetProjection(bson.M{
		"user_data." + m.Author.ID + ".score":        1,
		"user_data." + m.Author.ID + ".channel_data": bson.D{{Key: "$slice", Value: 10}},
	})
	findChannelData := database.ProcessedCollection.FindOne(ctx, filter, optionsFindChannelData)
	if findChannelData.Err() != nil {
		utils.NoStatsAviableForGuild(s, m)
		return
	}

	findChannelData.Decode(&guildObject)

	if guildObject.UserData[m.Author.ID].ChannelData == nil || len(guildObject.UserData[m.Author.ID].ChannelData) == 0 {
		utils.NoStatsAviableForYou(s, m)
		return
	}

	stats := MY_VOICE_STATS_MESSAGE_HEADING
	for _, value := range guildObject.UserData[m.Author.ID].ChannelData {
		stats += value.ChannelName + ": " + utils.FormatTime(value.Score) + "\n"
	}
	stats += "\nTotal time: " + utils.FormatTime(guildObject.UserData[m.Author.ID].Score) + "\n"

	_, _ = s.ChannelMessageSend(m.ChannelID, stats)
}
