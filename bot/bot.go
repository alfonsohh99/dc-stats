package bot

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
	"vc-stats/config"
	"vc-stats/constants"
	"vc-stats/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/procyon-projects/chrono"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var BotId string
var goBot *discordgo.Session
var dataCollection *mongo.Collection
var processedCollection *mongo.Collection
var fetchDataTask *chrono.ScheduledTask
var processDataTask *chrono.ScheduledTask
var ctx = context.TODO()

// DATA COLLECTION MODELS----
type Guild struct {
	ID      primitive.ObjectID `bson:"_id"`
	GuildID string             `bson:"guild_id"`
	Users   map[string]User    `bson:"users"`
}

type ChannelActivity struct {
	Score       uint64 `bson:"score"`
	ChannelName string `bson:"channel_name"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserID       string             `bson:"user_id"`
	UserName     string             `bson:"user_name"`
	UserActivity map[string]uint64  `bson:"user_activity"`
}

// PROCESSED COLLECTION MODELS
type UserScore struct {
	Username string `bson:"user_name"`
	Score    uint64 `bson:"score"`
}

type ChannelData struct {
	Score       uint64 `bson:"score"`
	ChannelName string `bson:"channel_name"`
}

type ProcessedGuild struct {
	ID       primitive.ObjectID       `bson:"_id"`
	GuildID  string                   `bson:"guild_id"`
	TopUsers []UserScore              `bson:"top_users"`
	UserData map[string][]ChannelData `bson:"user_data"`
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

	dataCollection = client.Database("discord").Collection("vc")
	processedCollection = client.Database("discord").Collection("vc-processed")

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

	taskScheduler := chrono.NewDefaultTaskScheduler()

	task, err := taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		gatherStats(goBot, ctx)
	}, constants.FetchDataInterval)
	fetchDataTask = &task

	if err == nil {
		log.Print("FetchDataTask has been scheduled successfully.  Fixed delay: ", constants.FetchDataInterval)
	}

	task, err = taskScheduler.ScheduleWithFixedDelay(func(ctx context.Context) {
		processStats(goBot, ctx)
	}, constants.ProcessDataInterval)
	processDataTask = &task

	if err == nil {
		log.Print("processDataTask has been scheduled successfully. Fixed delay: ", constants.ProcessDataInterval)
	}

	return

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == BotId && strings.Contains(m.Content, "pong_") {
		dateStr := strings.Split(m.Content, "_")[1]
		sentTime, _ := time.Parse(time.Layout, dateStr)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		_, _ = s.ChannelMessageSend(m.ChannelID, "ping "+m.Timestamp.Sub(sentTime).String())
		return
	}
	if m.Author.ID == BotId || !strings.Contains(m.Content, config.BotPrefix) {
		return
	}
	log.Println("message: ", m.Content)
	content := strings.Split(m.Content, config.BotPrefix)[1]

	switch content {
	case "ping":
		{
			_, _ = s.ChannelMessageSend(m.ChannelID, "pong_"+m.Timestamp.Format(time.Layout))
			break
		}
	case "myStats":
		{
			var guildObject Guild

			filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
			findGuild := dataCollection.FindOne(ctx, filter)
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

			for key, value := range guildObject.Users[m.Author.ID].UserActivity {
				channelNameFound := false
				for _, channel := range guildChannels {
					if channel.ID == key {
						stats += channel.Name + ", " + utils.FormatTime(value) + "\n"
						channelNameFound = true
						break
					}
				}
				if !channelNameFound {
					stats += "[id:" + key + "], " + utils.FormatTime(value) + "s\n"
				}

				totalTime += value
			}
			stats += "Total time: " + utils.FormatTime(totalTime) + "s\n"

			_, _ = s.ChannelMessageSend(m.ChannelID, stats)
			break
		}
	case "top":
		{
			var guildObject ProcessedGuild

			filter := bson.D{primitive.E{Key: "guild_id", Value: m.GuildID}}
			findProcessedGuild := processedCollection.FindOne(ctx, filter)

			if findProcessedGuild.Err() != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
				break
			}

			err := findProcessedGuild.Decode(&guildObject)

			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "No stats aviable for this guild")
				break
			}

			var response string
			response += ":beginner: SERVER RANKING :beginner: \n"
			for index, score := range guildObject.TopUsers {
				if index == 0 {
					response += ":trophy:"
				} else {
					response += "  " + strconv.Itoa(index+1) + "  "
				}
				response += " - " + score.Username + ": " + utils.FormatTime(score.Score) + "\n"
				if index == 10 {
					break
				}
			}
			response += "..."

			_, _ = s.ChannelMessageSend(m.ChannelID, response)
		}

	}
}

func gatherStats(goBot *discordgo.Session, ctx context.Context) {

	for _, guild := range goBot.State.Guilds {

		var guildObject Guild

		filter := bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findGuild := dataCollection.FindOne(ctx, filter)
		if findGuild.Err() != nil {
			//Guild doesnt exist
			log.Println("GUILD NOT PRESSENT", findGuild.Err().Error())
			newGuild := &Guild{
				ID:      primitive.NewObjectID(),
				GuildID: guild.ID,
				Users:   map[string]User{},
			}
			data, err := dataCollection.InsertOne(ctx, newGuild)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("CREATED GUILD, ", data)
			guildObject = *newGuild
		} else {

			findGuild.Decode(&guildObject)
		}

		members, err := goBot.GuildMembers(guild.ID, "", 1000)
		if err != nil {
			log.Println("ERROR GETTINGE GUILD MEMBERS")
			break
		}
		for _, member := range members {

			voiceState, err := goBot.State.VoiceState(guild.ID, member.User.ID)

			if err == nil {

				if guildObject.Users[member.User.ID].UserID == "" {
					log.Println("User not created")
					guildObject.Users[member.User.ID] = User{
						ID:           primitive.NewObjectID(),
						UserID:       member.User.ID,
						UserName:     member.Nick,
						UserActivity: map[string]uint64{voiceState.ChannelID: 10},
					}
					dataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
						{"$set", bson.D{{"users." + member.User.ID, guildObject.Users[member.User.ID]}}},
					})
				} else {
					currentValue := guildObject.Users[member.User.ID].UserActivity[voiceState.ChannelID]
					currentValue += 10
					dataCollection.UpdateByID(ctx, guildObject.ID, bson.D{
						{"$set", bson.D{{"users." + member.User.ID + ".user_activity." + voiceState.ChannelID, currentValue}}},
						{"$set", bson.D{{"users." + member.User.ID + ".user_name", member.Nick}}},
					})
				}

			}

		}
	}

}

func processStats(goBot *discordgo.Session, ctx context.Context) {

	for _, guild := range goBot.State.Guilds {

		var guildObject Guild

		filter := bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findGuild := dataCollection.FindOne(ctx, filter)
		if findGuild.Err() != nil {
			//Guild doesnt have data yet
			continue
		}

		findGuild.Decode(&guildObject)

		var scores []UserScore
		for _, user := range guildObject.Users {

			var total uint64
			for _, value := range user.UserActivity {
				total += value
			}
			scores = append(scores, UserScore{Username: user.UserName, Score: total})

		}
		sort.SliceStable(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})
		var processedGuildObject ProcessedGuild

		filter = bson.D{primitive.E{Key: "guild_id", Value: guild.ID}}
		findProcessedGuild := processedCollection.FindOne(ctx, filter)
		if findProcessedGuild.Err() != nil {
			//Guild doesnt exist
			log.Println("PROCESSED GUILD NOT PRESSENT", findProcessedGuild.Err().Error())
			newGuild := ProcessedGuild{
				ID:       primitive.NewObjectID(),
				GuildID:  guild.ID,
				TopUsers: scores,
			}
			data, err := processedCollection.InsertOne(ctx, newGuild)
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("CREATED PROCESSED GUILD, ", data)
			processedGuildObject = newGuild
		} else {
			findProcessedGuild.Decode(&processedGuildObject)
			processedCollection.UpdateByID(ctx, processedGuildObject.ID, bson.D{
				{"$set", bson.D{{"top_users", scores}}},
			})
		}

		return

		var response string
		response += ":beginner: SERVER RANKING :beginner: \n"
		for index, score := range scores {
			if index == 0 {
				response += ":trophy:"
			} else {
				response += "  " + strconv.Itoa(index+1) + "  "
			}
			response += " - " + score.Username + ": " + utils.FormatTime(score.Score) + "\n"
			if index == 10 {
				break
			}
		}
		response += "..."

	}

}
