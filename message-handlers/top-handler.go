package messagehandlers

import (
	"context"
	"strconv"
	"vc-stats/database"
	"vc-stats/model"
	"vc-stats/utils"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Top(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {
	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	findProcessedGuild := database.ProcessedCollection.FindOne(ctx, filter)

	if findProcessedGuild.Err() != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	err := findProcessedGuild.Decode(&guildObject)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	var response string
	response += ":beginner: SERVER RANKING :beginner: \n"
	for index, score := range guildObject.TopUsers {
		if index == 0 {
			response += ":trophy:"
		} else {
			response += "  " + strconv.Itoa(index+1) + "  "
		}
		response += " - " + score.Username + ": " + utils.FormatTime(score.Score) + "\n"
		if index == 10 {
			return
		}
	}
	response += "..."

	_, _ = s.ChannelMessageSend(m.ChannelID, response)
}
