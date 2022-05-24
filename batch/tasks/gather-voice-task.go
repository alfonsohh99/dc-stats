package tasks

import (
	"context"
	"dc-stats/database"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func GatherVoiceStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {

		guildObject, err := database.FindOrCreateDataGuild(guild.ID, ctx)
		if err != nil {
			log.Println("Error finding/creating guild", err)
			continue
		}

		lastId := ""
		for {
			members, err := goBot.GuildMembers(guild.ID, lastId, 1000)
			if err != nil || len(members) == 0 {
				break
			}
			lastId = members[len(members)-1].User.ID
			for _, member := range members {

				nickName := member.Nick
				if nickName == "" {
					nickName = member.User.Username
				}

				voiceState, err := goBot.State.VoiceState(guild.ID, member.User.ID)

				if err == nil {
					database.SaveOrUpdateDataGuildVoiceState(guildObject, member.User.ID, voiceState.ChannelID, nickName, ctx)
				}

			}
		}
	}

}
