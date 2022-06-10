package bot

import (
	"context"
	"dc-stats/bot/message-handlers"
	"dc-stats/config"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session
var ctx context.Context

func Start(context context.Context) (bot *discordgo.Session) {

	ctx = context
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

	err = goBot.Open()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Bot is running !")
	return goBot

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

	content := strings.Split(m.Content, config.BotPrefix)[1]

	switch content {
	case "ping":
		{
			go messagehandlers.Ping(s, m)
			break
		}
	case "myVoice":
		{
			go messagehandlers.MyVoiceStats(s, m, ctx)
			break
		}
	case "myMessage":
		{
			go messagehandlers.MyMessageStats(s, m, ctx)
			break
		}
	case "topVoice":
		{
			go messagehandlers.TopVoice(s, m, ctx)
			break
		}
	case "topMessage":
		{
			go messagehandlers.TopMessage(s, m, ctx)
			break
		}

	}
}
