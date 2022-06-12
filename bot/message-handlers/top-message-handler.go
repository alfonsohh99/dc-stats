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

const TOP_MESSAGE_MESSAGE_HEADING = "**:newspaper: SERVER CHAT MESSAGE RANKING :newspaper:**\n\n"

func TopMessage(s *discordgo.Session, m *discordgo.MessageCreate, ctx context.Context) {

	var guildObject processedModel.Guild

	filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
	opts := options.FindOne().SetProjection(bson.M{
		"_id":               1,
		"top_message_users": bson.D{{Key: "$slice", Value: 10}},
		// We need to project another field other than top users so that it excludes the rest and $slice works properly
	})
	findTopMessageUsers := database.ProcessedCollection.FindOne(ctx, filter, opts)

	if findTopMessageUsers.Err() != nil {
		utils.NoStatsAviableForGuild(s, m, findTopMessageUsers.Err())
		return
	}

	errProcess := findTopMessageUsers.Decode(&guildObject)

	if errProcess != nil || len(guildObject.TopMessageUsers) == 0 {
		if errProcess != nil {
			utils.NoStatsAviableForGuild(s, m, errProcess)
		}
		return
	}

	response := TOP_MESSAGE_MESSAGE_HEADING
	for index, score := range guildObject.TopMessageUsers {
		if index == 0 {
			response += TROHPY_EMOJI
		} else {
			response += "  " + strconv.Itoa(index+1) + "  "
		}
		response += " - **" + score.Username + "**: " + strconv.FormatUint(score.Score, 10) + " messages \n"
	}

	_, _ = s.ChannelMessageSend(m.ChannelID, response)
}
