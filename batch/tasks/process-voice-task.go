package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"log"
	"sort"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func ProcessVoiceStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {

		guildId := guild.ID

		guildObject, err := database.FindDataGuild(ctx, guildId)
		if err != nil {
			log.Println("Cannot find guild to process")
			continue
		}

		scores := []model.UserScore{}
		userData := map[string]model.ProcessedUser{}
		for _, user := range guildObject.Users {
			channelData := []model.ChannelData{}
			var totalScore uint64
			for channelId, value := range user.UserVoiceActivity {
				channelName := guildObject.ChannelMarks[channelId].Name
				if channelName == "" {
					channelName = "[" + channelId + "], "
				}
				channelData = append(channelData, model.ChannelData{ChannelName: channelName, Score: value})
				totalScore += value

			}
			if totalScore != 0 {
				scores = append(scores, model.UserScore{Username: user.UserName, Score: totalScore})
				sort.SliceStable(channelData, func(i, j int) bool {
					return channelData[i].Score > channelData[j].Score
				})
				userData[user.UserID] = model.ProcessedUser{Score: totalScore, ChannelData: channelData}
			}

		}

		sort.SliceStable(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})

		database.SaveOrUpdateProcessedGuild(guildId, scores, userData, ctx)

	}

}
