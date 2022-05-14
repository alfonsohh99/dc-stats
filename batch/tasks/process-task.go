package tasks

import (
	"context"
	"log"
	"sort"
	"sync"
	"vc-stats/database"
	"vc-stats/model"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProcessStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {

		guildChannels, err := goBot.GuildChannels(guild.ID)
		if err != nil {
			// Error retrieving guild channels
			log.Println("CANNOT ACCESS GUILD CHANNELS")
			continue
		}

		var guildObject model.Guild

		filter := bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findGuild := database.DataCollection.FindOne(ctx, filter)
		if findGuild.Err() != nil {
			log.Println("Cannot find guild to process")
			//Guild doesnt have data yet
			continue
		}

		findGuild.Decode(&guildObject)

		scores := []model.UserScore{}
		userData := map[string]model.ProcessedUser{}
		for _, user := range guildObject.Users {
			//CALCULATING TOTAL SCORE
			var total uint64
			for _, value := range user.UserActivity {
				total += value
			}
			scores = append(scores, model.UserScore{Username: user.UserName, Score: total})

			// CALCULATING CHANNEL DATA
			channelData := []model.ChannelData{}

			for key, value := range guildObject.Users[user.UserID].UserActivity {
				channelNameFound := false
				for _, channel := range guildChannels {
					if channel.ID == key {
						channelData = append(channelData, model.ChannelData{ChannelName: channel.Name, Score: value})
						channelNameFound = true
						break
					}
				}
				if !channelNameFound {
					channelData = append(channelData, model.ChannelData{ChannelName: "[" + key + "], ", Score: value})
				}

			}
			sort.SliceStable(channelData, func(i, j int) bool {
				return channelData[i].Score > channelData[j].Score
			})
			userData[user.UserID] = model.ProcessedUser{Score: total, ChannelData: channelData}

		}
		sort.SliceStable(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})
		var processedGuildObject model.ProcessedGuild

		filter = bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findProcessedGuild := database.ProcessedCollection.FindOne(ctx, filter)
		if findProcessedGuild.Err() != nil {
			//Guild doesnt exist
			log.Println("PROCESSED GUILD NOT PRESSENT", findProcessedGuild.Err().Error())
			newGuild := model.ProcessedGuild{
				ID:       primitive.NewObjectID(),
				GuildID:  guild.ID,
				TopUsers: scores,
				UserData: userData,
			}
			data, err := database.ProcessedCollection.InsertOne(ctx, newGuild)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("CREATED PROCESSED GUILD, ", data)
			processedGuildObject = newGuild
		} else {
			findProcessedGuild.Decode(&processedGuildObject)
			database.ProcessedCollection.UpdateByID(ctx, processedGuildObject.ID, bson.D{
				{"$set", bson.D{{"top_users", scores}}},
				{"$set", bson.D{{"user_data", userData}}},
			})
		}

	}

}
