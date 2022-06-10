package messagehandlers

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"dc-stats/utils"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MY_MESSAGE_STATS_MESSAGE_HEADING = "**:newspaper: YOUR TOP ACTIVE CHANNELS :newspaper:**\n\n"

func MyMessageStats(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	optionsFindChannelData := options.FindOne().SetProjection(bson.M{
		"user_message_data." + m.Author.ID + ".score":        1,
		"user_message_data." + m.Author.ID + ".channel_data": bson.D{{Key: "$slice", Value: 10}},
	})
	findChannelData := database.ProcessedCollection.FindOne(ctx, filter, optionsFindChannelData)
	if findChannelData.Err() != nil {
		utils.NoStatsAviableForGuild(s, m, findChannelData.Err())
		return
	}

	findChannelData.Decode(&guildObject)

	if guildObject.UserMessageData[m.Author.ID].ChannelData == nil || len(guildObject.UserMessageData[m.Author.ID].ChannelData) == 0 {
		utils.NoStatsAviableForYou(s, m)
		return
	}

	stats := MY_MESSAGE_STATS_MESSAGE_HEADING
	for index, value := range guildObject.UserMessageData[m.Author.ID].ChannelData {
		if index == 0 {
			stats += TROHPY_EMOJI
		} else {
			stats += "  " + strconv.Itoa(index+1) + "  "
		}
		stats += " - **" + value.ChannelName + "**: " + strconv.FormatUint(value.Score, 10) + " messages \n"
	}
	stats += "\nTotal messages: " + strconv.FormatUint(guildObject.UserMessageData[m.Author.ID].Score, 10) + " messages \n"

	_, _ = s.ChannelMessageSend(m.ChannelID, stats)
}
