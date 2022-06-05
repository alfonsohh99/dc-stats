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

const TOP_VOICE_MESSAGE_HEADING = ":beginner: SERVER VOICE CHAT RANKING :beginner: \n\n"
const TROHPY_EMOJI = ":trophy:"

func TopVoice(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject model.ProcessedGuild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":       1,
		"top_users": bson.D{{Key: "$slice", Value: 10}},
		// We need to project another field other than top users so that it excludes the rest and $slice works properly
	})
	findTopUsers := database.ProcessedCollection.FindOne(ctx, filter, opts)

	if findTopUsers.Err() != nil {
		utils.NoStatsAviableForGuild(s, m)
		return
	}

	errProcess := findTopUsers.Decode(&guildObject)

	if errProcess != nil || len(guildObject.TopUsers) == 0 {
		utils.NoStatsAviableForGuild(s, m)
		return
	}

	response := TOP_VOICE_MESSAGE_HEADING
	for index, score := range guildObject.TopUsers {
		if index == 0 {
			response += TROHPY_EMOJI
		} else {
			response += "  " + strconv.Itoa(index+1) + "  "
		}
		response += " - " + score.Username + ": " + utils.FormatTime(score.Score) + "\n"
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, response)
}
