package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

func GatherMessageStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {

		guildObject, err := database.FindOrCreateDataGuild(guild.ID, ctx)
		if err != nil {
			log.Println("Error finding/creating guild", err)
			continue
		}

		for _, channel := range guild.Channels {
			if channel.Type != 0 {
				continue
			}
			channelMark := guildObject.ChannelMarks[channel.ID]
			if channelMark.BeforeId == "" && channelMark.AfterId != "" {
				messages, err := goBot.ChannelMessages(channel.ID, 100, "", channelMark.AfterId, "")
				if err != nil || len(messages) == 0 {
					if err != nil {
						log.Println("Error getting channel mesasges, forward, ", channel.ID, err)
					}
					continue
				}
				processMessages(messages, channel.ID, &guildObject, ctx)
				channelMark.AfterId = messages[0].ID
				channelMark.TotalMessages += uint64(len(messages))

			} else {

				messages, err := goBot.ChannelMessages(channel.ID, 100, channelMark.BeforeId, "", "")

				if err != nil {
					log.Println("Error getting channel mesasges, backwards, ", channel.ID, err)
					continue
				}

				if len(messages) == 0 {
					if channelMark.BeforeId != "" {
						log.Println("FINISHED ANALYZING BACKWARDS: ", channel.Name)
						channelMark.BeforeId = ""
					}
					guildObject.ChannelMarks[channel.ID] = channelMark
					continue
				}

				processMessages(messages, channel.ID, &guildObject, ctx)
				channelMark.TotalMessages += uint64(len(messages))
				if channelMark.BeforeId == "" {
					channelMark.AfterId = messages[0].ID
				}
				channelMark.BeforeId = messages[len(messages)-1].ID
			}
			guildObject.ChannelMarks[channel.ID] = channelMark

		}
		database.UpdateDataGuildUsersAndChannelMarks(guildObject, ctx)
	}

}

func processMessages(messagess []*discordgo.Message, channelId string, guildObject *model.Guild, ctx context.Context) {
	for _, message := range messagess {
		if guildObject.Users[message.Author.ID].UserID == "" {
			log.Println("User not created")
			user := model.CreateDataUser(message.Author.ID)
			user.UserMessageActivity[channelId] = 1
			guildObject.Users[message.Author.ID] = user
		} else {
			currentValue := guildObject.Users[message.Author.ID].UserMessageActivity[channelId] + 1
			guildObject.Users[message.Author.ID].UserMessageActivity[channelId] = currentValue
		}
	}

}
