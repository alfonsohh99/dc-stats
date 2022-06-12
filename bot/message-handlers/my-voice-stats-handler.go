package messagehandlers

import (
	"context"
	"dc-stats/database"
	"dc-stats/model/processed"
	"dc-stats/utils"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MY_VOICE_STATS_MESSAGE_HEADING = "**:speaking_head: YOUR TOP ACTIVE VOICE CHANNELS :speaking_head:**\n\n"

func MyVoiceStats(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject processedModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	optionsFindChannelData := options.FindOne().SetProjection(bson.M{
		"user_data." + m.Author.ID + ".score":        1,
		"user_data." + m.Author.ID + ".channel_data": bson.D{{Key: "$slice", Value: 10}},
	})
	findChannelData := database.ProcessedCollection.FindOne(ctx, filter, optionsFindChannelData)
	if findChannelData.Err() != nil {
		utils.NoStatsAviableForGuild(s, m, findChannelData.Err())
		return
	}

	findChannelData.Decode(&guildObject)

	if guildObject.UserData[m.Author.ID].ChannelData == nil || len(guildObject.UserData[m.Author.ID].ChannelData) == 0 {
		utils.NoStatsAviableForYou(s, m)
		return
	}

	stats := MY_VOICE_STATS_MESSAGE_HEADING
	for index, value := range guildObject.UserData[m.Author.ID].ChannelData {
		if index == 0 {
			stats += TROHPY_EMOJI
		} else {
			stats += "  " + strconv.Itoa(index+1) + "  "
		}
		stats += " - **" + value.ChannelName + "**: " + utils.FormatTime(value.Score) + "\n"
	}
	stats += "\nTotal time: " + utils.FormatTime(guildObject.UserData[m.Author.ID].Score) + "\n"

	_, _ = s.ChannelMessageSend(m.ChannelID, stats)
}
