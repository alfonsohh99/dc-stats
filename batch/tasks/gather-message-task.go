package tasks

import (
	"context"
	"dc-stats/database"
	"dc-stats/model"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
			if channel.Type != 0 { // && channel.Type != 2
				continue
			}
			channelMark := guildObject.ChannelMarks[channel.ID]
			if channelMark.BeforeId == "" && channelMark.AfterId != "" {
				messages, err := goBot.ChannelMessages(channel.ID, 100, "", channelMark.AfterId, "")
				if err != nil {
					log.Println("Error getting channel mesasges, forward, ", channel.ID)
					continue
				}
				if len(messages) == 0 {
					continue
				}
				processMessages(messages, channel.ID, &guildObject, ctx)
				channelMark.AfterId = messages[0].ID
				channelMark.TotalMessages += uint64(len(messages))
				guildObject.ChannelMarks[channel.ID] = channelMark

			} else {

				messages, err := goBot.ChannelMessages(channel.ID, 100, channelMark.BeforeId, "", "")

				if err != nil {
					log.Println("Error getting channel mesasges, backwards, ", channel.ID)
					continue
				}

				if len(messages) == 0 {
					if channelMark.BeforeId != "" {
						log.Println("FINISHED ANALYZING BACKWARDS: ", channel.Name)
						channelMark.BeforeId = ""
						guildObject.ChannelMarks[channel.ID] = channelMark
					}
					continue
				}

				processMessages(messages, channel.ID, &guildObject, ctx)
				channelMark.TotalMessages += uint64(len(messages))
				if channelMark.BeforeId == "" {
					channelMark.AfterId = messages[0].ID
				}
				channelMark.BeforeId = messages[len(messages)-1].ID
				guildObject.ChannelMarks[channel.ID] = channelMark
			}

		}
		database.DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
			{"$set", bson.D{{"channel_marks", guildObject.ChannelMarks}}},
			{"$set", bson.D{{"users", guildObject.Users}}},
		})

	}

}

func processMessages(messagess []*discordgo.Message, channelId string, guildObject *model.Guild, ctx context.Context) {
	for _, message := range messagess {
		log.Println("MESSAGE")
		nickName := message.Author.Username
		// TODO, EN LA CACHE QUE HAREMOS DE ID -> NICKNAME DE USAURIOS MIRAR EL NICK
		if guildObject.Users[message.Author.ID].UserID == "" {
			log.Println("User not created")
			user := model.User{
				ID:                  primitive.NewObjectID(),
				UserID:              message.Author.ID,
				UserName:            nickName,
				UserVoiceActivity:   map[string]uint64{},
				UserMessageActivity: map[string]uint64{channelId: 1},
			}
			guildObject.Users[message.Author.ID] = user
		} else {
			currentValue := guildObject.Users[message.Author.ID].UserMessageActivity[channelId] + 1
			guildObject.Users[message.Author.ID].UserMessageActivity[channelId] = currentValue
		}
	}

}
