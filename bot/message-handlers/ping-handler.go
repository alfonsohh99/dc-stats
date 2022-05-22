package messagehandlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, "pong_"+m.Timestamp.Format(time.Layout))
}
