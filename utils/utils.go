package utils

import (
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
		return strconv.FormatUint(timeSeconds, 10) + " seconds "
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
