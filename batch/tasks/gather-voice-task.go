package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func GatherVoiceStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {
		guildId := guild.ID
		guildObject, err := database.FindOrCreateDataGuild(guildId, ctx)
		if err != nil {
			log.Println("Error finding/creating guild", err)
			continue
		}

		lastId := ""
		for {
			// TODO: Guilds with more than 1K members may experience inconsistent time measuring
			// Proposal: Do not iterate over members, iterate over voice channels and check if there are users inside
			members, err := goBot.GuildMembers(guildId, lastId, 1000)
			if err != nil || len(members) == 0 {
				break
			}
			lastId = members[len(members)-1].User.ID
			for _, member := range members {

				nickName := member.Nick
				if nickName == "" {
					nickName = member.User.Username
				}
				userId := member.User.ID

				voiceState, err := goBot.State.VoiceState(guildId, userId)

				if err == nil {
					channelId := voiceState.ChannelID
					savedUser := guildObject.Users[userId]
					if savedUser.UserID == "" {
						newUser := model.CreateDataUser(userId, nickName)
						newUser.UserVoiceActivity[channelId] = 10
						guildObject.Users[userId] = newUser
					} else {
						currentValue := savedUser.UserVoiceActivity[channelId] + 10
						savedUser.UserVoiceActivity[channelId] = currentValue
						savedUser.UserName = nickName
					}
				}

			}
			database.UpdateDataGuildUsers(guildObject, ctx)
		}
	}

}
