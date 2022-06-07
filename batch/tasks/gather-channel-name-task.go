package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func GatherChannelNameStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {
		guildId := guild.ID
		guildObject, err := database.FindOrCreateDataGuild(guildId, ctx)
		if err != nil {
			log.Println("Error finding/creating guild", err)
			continue
		}
		for _, channel := range guild.Channels {
			mark, exists := guildObject.ChannelMarks[channel.ID]
			if !exists {
				newMark := model.MessageMark{
					BeforeId:      "",
					AfterId:       "",
					TotalMessages: 0,
					Name:          channel.Name,
				}
				guildObject.ChannelMarks[channel.ID] = newMark
			} else {
				mark.Name = channel.Name
				guildObject.ChannelMarks[channel.ID] = mark
			}
		}
		database.UpdateDataGuildChannelMarks(guildObject, ctx)
	}

}
