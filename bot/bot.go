package bot

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"vc-stats/config"
	"vc-stats/constants"
)

var BotId string
var goBot *discordgo.Session
var collection *mongo.Collection
var ctx = context.TODO()

type Guild struct {
	ID      primitive.ObjectID `bson:"_id"`
	GuildID string             `bson:"guild_id"`
	Users   map[string]User    `bson:"users"`
}

type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	UserID          string             `bson:"user_id"`
	ChannelActivity map[string]uint64  `bson:"channel_activity"`
}

func Start() {

	clientOptions := options.Client().ApplyURI("mongodb://" + config.DatabaseEndpoint + ":" + config.DatabasePort + "/").SetAuth(options.Credential{
		Username: config.DatabaseUser,
		Password: config.DatabasePassword})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("discord").Collection("vc")

	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		log.Fatal(err)
	}

	u, err := goBot.User("@me")

	if err != nil {
		log.Fatal(err)
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = goBot.Open()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running !")

	ticker := time.NewTicker(constants.FetchDataInterval)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			gatherStats(goBot)
		case <-quit:
			ticker.Stop()
			return
		}
	}

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotId || !strings.Contains(m.Content, config.BotPrefix) {
		return
	}
	log.Println("message: ", m.Content)
	content := strings.Split(m.Content, config.BotPrefix)[1]

	switch content {
	case "ping":
		{
			_, _ = s.ChannelMessageSend(m.ChannelID, "pong pong")
			break
		}
	case "myStats":
		{
			var guildObject Guild

			filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
			findGuild := collection.FindOne(ctx, filter)
			guildChannels, err := s.GuildChannels(m.GuildID)

			if findGuild.Err() != nil || err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
				break
			}

			findGuild.Decode(&guildObject)

			if guildObject.Users[m.Author.ID].UserID == "" {
				_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for you :(")
				break
			}

			var stats string
			stats += ":beginner: YOUR VOICE CHAT STATS :beginner:\n"

			var totalTime uint64

			for key, value := range guildObject.Users[m.Author.ID].ChannelActivity {
				channelNameFound := false
				for _, channel := range guildChannels {
					if channel.ID == key {
						stats += channel.Name + ", " + formatTime(value) + "\n"
						channelNameFound = true
						break
					}
				}
				if !channelNameFound {
					stats += "[id:" + key + "], " + formatTime(value) + "s\n"
				}

				totalTime += value
			}
			stats += "Total time: " + formatTime(totalTime) + "s\n"

			_, _ = s.ChannelMessageSend(m.ChannelID, stats)
			break
		}
	case "top":
		{
			var guildObject Guild

			filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
			findGuild := collection.FindOne(ctx, filter)

			if findGuild.Err() != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
				break
			}

			findGuild.Decode(&guildObject)

			type UserScore struct {
				Username string
				Score    uint64
			}
			var scores []UserScore
			for _, user := range guildObject.Users {

				guildMembers, err := s.GuildMembers(guildObject.GuildID, "", 1000)

				if err != nil {
					log.Fatal(err)
				}

				var username string
				for _, guildUser := range guildMembers {
					if guildUser.User.ID == user.UserID {
						username = guildUser.Nick
						break
					}
				}
				var total uint64
				for _, value := range user.ChannelActivity {
					total += value
				}
				scores = append(scores, UserScore{Username: username, Score: total})
			}
			sort.SliceStable(scores, func(i, j int) bool {
				return scores[i].Score > scores[j].Score
			})
			var response string
			response += ":beginner: SERVER RANKING :beginner: \n"
			for index, score := range scores {
				if index == 0 {
					response += ":trophy:"
				} else {
					response += "  " + strconv.Itoa(index+1) + "  "
				}
				response += " - " + score.Username + ": " + formatTime(score.Score) + "\n"
				if index == 10 {
					break
				}
			}
			response += "..."

			_, _ = s.ChannelMessageSend(m.ChannelID, response)

		}

	}
}

func gatherStats(goBot *discordgo.Session) {

	for _, guild := range goBot.State.Guilds {

		var guildObject Guild

		filter := bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findGuild := collection.FindOne(ctx, filter)
		if findGuild.Err() != nil {
			//Guild doesnt exist
			log.Println("GUILD NOT PRESSENT", findGuild.Err().Error())
			newGuild := &Guild{
				ID:      primitive.NewObjectID(),
				GuildID: guild.ID,
				Users:   map[string]User{},
			}
			data, err := collection.InsertOne(ctx, newGuild)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("CREATED GUILD, ", data)
			guildObject = *newGuild
		} else {

			findGuild.Decode(&guildObject)
		}

		members, err := goBot.GuildMembers(guild.ID, "", 1000)
		if err != nil {
			log.Fatal("ERROR GETTINGE GUILD MEMBERS")
		}
		for _, member := range members {

			voiceState, err := goBot.State.VoiceState(guild.ID, member.User.ID)

			if err == nil {

				if guildObject.Users[member.User.ID].UserID == "" {
					log.Println("User not created")
					guildObject.Users[member.User.ID] = User{
						ID:              primitive.NewObjectID(),
						UserID:          member.User.ID,
						ChannelActivity: map[string]uint64{voiceState.ChannelID: 10},
					}
					collection.UpdateByID(ctx, guildObject.ID, bson.D{
						{"$set", bson.D{{"users." + member.User.ID, guildObject.Users[member.User.ID]}}},
					})
				} else {
					guildObject.Users[member.User.ID].ChannelActivity[voiceState.ChannelID] += 10
					collection.UpdateByID(ctx, guildObject.ID, bson.D{
						{"$set", bson.D{{"users." + member.User.ID + ".channel_activity." + voiceState.ChannelID, guildObject.Users[member.User.ID].ChannelActivity[voiceState.ChannelID]}}},
					})
				}

			}

		}
	}

}

func formatTime(timeSeconds uint64) string {

	var res string
	var seconds uint64
	var minutes uint64
	var hours uint64
	var days uint64
	var years uint64
	if timeSeconds < 60 {
		return strconv.FormatUint(timeSeconds, 10)
	}
	seconds = timeSeconds % 60
	minutes = ((timeSeconds - seconds) / 60) % 60
	hours = ((timeSeconds - seconds - minutes*60) / 3600) % 24
	days = ((timeSeconds - seconds - minutes*60 - hours*3600) / 86400) % 365
	years = (timeSeconds - seconds - minutes*60 - hours*3600 - days*86400) / 31536000

	if years > 0 {
		res += strconv.FormatUint(years, 10) + " years "
	}

	if days > 0 {
		res += strconv.FormatUint(days, 10) + " days "
	}

	if hours > 0 {
		res += strconv.FormatUint(hours, 10) + " hours "
	}

	if minutes > 0 {
		res += strconv.FormatUint(minutes, 10) + " minutes "
	}

	if seconds > 0 {
		res += strconv.FormatUint(seconds, 10) + " seconds "
	}

	return res
}
