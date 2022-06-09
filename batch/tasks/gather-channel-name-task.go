package tasks

import (
	"context"
	"dc-stats/database"
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
			guildObject.ChannelNameMap[channel.ID] = channel.Name
		}
		database.UpdateDataGuildChannelNameMap(guildObject, ctx)
	}

}
