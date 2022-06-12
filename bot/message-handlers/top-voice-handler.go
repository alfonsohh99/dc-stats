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

const TOP_VOICE_MESSAGE_HEADING = "**:speaking_head: SERVER VOICE CHAT RANKING :speaking_head:**\n\n"
const TROHPY_EMOJI = ":trophy:"

func TopVoice(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject processedModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":       1,
		"top_users": bson.D{{Key: "$slice", Value: 10}},
		// We need to project another field other than top users so that it excludes the rest and $slice works properly
	})
	findTopUsers := database.ProcessedCollection.FindOne(ctx, filter, opts)

	if findTopUsers.Err() != nil {
		utils.NoStatsAviableForGuild(s, m, findTopUsers.Err())
		return
	}

	errProcess := findTopUsers.Decode(&guildObject)

	if errProcess != nil || len(guildObject.TopUsers) == 0 {
		if errProcess != nil {
			utils.NoStatsAviableForGuild(s, m, errProcess)
		}
		return
	}

	response := TOP_VOICE_MESSAGE_HEADING
	for index, score := range guildObject.TopUsers {
		if index == 0 {
			response += TROHPY_EMOJI
		} else {
			response += "  " + strconv.Itoa(index+1) + "  "
		}
		response += " - **" + score.Username + "**: " + utils.FormatTime(score.Score) + "\n"
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, response)
}
