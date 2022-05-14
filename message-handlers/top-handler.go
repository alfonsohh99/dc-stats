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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Top(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {
	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":       1,
		"top_users": bson.D{{Key: "$slice", Value: 10}},
	})
	findTopUsers := database.ProcessedCollection.FindOne(ctx, filter, opts)

	if findTopUsers.Err() != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	errProcess := findTopUsers.Decode(&guildObject)

	if errProcess != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
		return
	}

	var response string
	response += ":beginner: SERVER RANKING :beginner: \n\n"
	for index, score := range guildObject.TopUsers {
		if index == 0 {
			response += ":trophy:"
		} else {
			response += "  " + strconv.Itoa(index+1) + "  "
		}
		response += " - " + score.Username + ": " + utils.FormatTime(score.Score) + "\n"
	}
	response += "..."

	_, _ = s.ChannelMessageSend(m.ChannelID, response)
}
