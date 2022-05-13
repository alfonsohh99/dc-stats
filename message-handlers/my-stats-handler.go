package messagehandlers

import (
	"context"
	"vc-stats/database"
	"vc-stats/model"
	"vc-stats/utils"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MyStats(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	findGuild := database.ProcessedCollection.FindOne(ctx, filter)
	if findGuild.Err() != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	findGuild.Decode(&guildObject)

	if guildObject.UserData[m.Author.ID].ChannelData == nil || len(guildObject.UserData[m.Author.ID].ChannelData) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for you :(")
		return
	}

	var stats string
	stats += ":beginner: YOUR VOICE CHAT STATS :beginner:\n"

	for _, value := range guildObject.UserData[m.Author.ID].ChannelData {
		stats += value.ChannelName + ": " + utils.FormatTime(value.Score) + "\n"
	}
	stats += "Total time: " + utils.FormatTime(guildObject.UserData[m.Author.ID].Score) + "s\n"

	_, _ = s.ChannelMessageSend(m.ChannelID, stats)
}
