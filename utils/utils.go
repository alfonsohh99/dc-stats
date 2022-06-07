package utils

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func FormatTime(timeSeconds uint64) string {

	var res string
	var seconds uint64
	var minutes uint64
	var hours uint64
	var days uint64
	var years uint64
	if timeSeconds < 60 {
		return strconv.FormatUint(timeSeconds, 10) + "s "
	}
	seconds = timeSeconds % 60
	minutes = ((timeSeconds - seconds) / 60) % 60
	hours = ((timeSeconds - seconds - minutes*60) / 3600) % 24
	days = ((timeSeconds - seconds - minutes*60 - hours*3600) / 86400) % 365
	years = (timeSeconds - seconds - minutes*60 - hours*3600 - days*86400) / 31536000

	if years > 0 {
		res += strconv.FormatUint(years, 10) + "y "
	}

	if days > 0 {
		res += strconv.FormatUint(days, 10) + "d "
	}

	if hours > 0 {
		res += strconv.FormatUint(hours, 10) + "h "
	}

	if minutes > 0 {
		res += strconv.FormatUint(minutes, 10) + "m "
	}

	if seconds > 0 {
		res += strconv.FormatUint(seconds, 10) + "s "
	}

	return res
}

func NoStatsAviableForGuild(s *discordgo.Session, message *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(message.ChannelID, "No stats aviable for this guild")
}

func NoStatsAviableForYou(s *discordgo.Session, message *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(message.ChannelID, "No stats aviable for you")
}

func GetUserNickName(member discordgo.Member) string {
	nickName := member.Nick
	if nickName == "" {
		nickName = member.User.Username
	}
	return nickName
}
