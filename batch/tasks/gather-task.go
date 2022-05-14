package tasks

import (
	"context"
	"log"
	"sync"
	"vc-stats/database"
	"vc-stats/model"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GatherStats(goBot *discordgo.Session, ctx context.Context, wait *sync.WaitGroup) {
	defer wait.Done()

	for _, guild := range goBot.State.Guilds {

		var guildObject model.Guild

		filter := bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findGuild := database.DataCollection.FindOne(ctx, filter)
		if findGuild.Err() != nil {
			//Guild doesnt exist
			log.Println("GUILD NOT PRESSENT", findGuild.Err().Error())
			newGuild := &model.Guild{
				ID:      primitive.NewObjectID(),
				GuildID: guild.ID,
				Users:   map[string]model.User{},
			}
			data, err := database.DataCollection.InsertOne(ctx, newGuild)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("CREATED GUILD, ", data)
			guildObject = *newGuild
		} else {

			findGuild.Decode(&guildObject)
		}
		lastId := ""
		for {
			members, err := goBot.GuildMembers(guild.ID, lastId, 1000)
			if err != nil || len(members) == 0 {
				break
			}
			lastId = members[len(members)-1].User.ID
			for _, member := range members {

				voiceState, err := goBot.State.VoiceState(guild.ID, member.User.ID)

				if err == nil {

					if guildObject.Users[member.User.ID].UserID == "" {
						log.Println("User not created")
						guildObject.Users[member.User.ID] = model.User{
							ID:           primitive.NewObjectID(),
							UserID:       member.User.ID,
							UserName:     member.Nick,
							UserActivity: map[string]uint64{voiceState.ChannelID: 10},
						}
						database.DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
							{"$set", bson.D{{"users." + member.User.ID, guildObject.Users[member.User.ID]}}},
						})
					} else {
						currentValue := guildObject.Users[member.User.ID].UserActivity[voiceState.ChannelID]
						currentValue += 10
						database.DataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
							{"$set", bson.D{{"users." + member.User.ID + ".user_activity." + voiceState.ChannelID, currentValue}}},
							{"$set", bson.D{{"users." + member.User.ID + ".user_name", member.Nick}}},
						})
					}

				}

			}
		}
	}

}
