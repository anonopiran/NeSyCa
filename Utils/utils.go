package utils

import (
	"fmt"
	"math"
	"strconv"
)

func ByteCountIEC(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func Plural(count int, singular string) (result string) {
	if (count == 1) || (count == 0) {
		result = strconv.Itoa(count) + " " + singular + " "
	} else {
		result = strconv.Itoa(count) + " " + singular + "s "
	}
	return
}

func SecondsToHuman(input int) (result string) {
	years := math.Floor(float64(input) / 60 / 60 / 24 / 7 / 30 / 12)
	seconds := input % (60 * 60 * 24 * 7 * 30 * 12)
	months := math.Floor(float64(seconds) / 60 / 60 / 24 / 7 / 30)
	seconds = input % (60 * 60 * 24 * 7 * 30)
	weeks := math.Floor(float64(seconds) / 60 / 60 / 24 / 7)
	seconds = input % (60 * 60 * 24 * 7)
	days := math.Floor(float64(seconds) / 60 / 60 / 24)
	seconds = input % (60 * 60 * 24)
	hours := math.Floor(float64(seconds) / 60 / 60)
	seconds = input % (60 * 60)
	minutes := math.Floor(float64(seconds) / 60)
	seconds = input % 60

	if years > 0 {
		result = Plural(int(years), "year") + Plural(int(months), "month") + Plural(int(weeks), "week") + Plural(int(days), "day") + Plural(int(hours), "hour") + Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else if months > 0 {
		result = Plural(int(months), "month") + Plural(int(weeks), "week") + Plural(int(days), "day") + Plural(int(hours), "hour") + Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else if weeks > 0 {
		result = Plural(int(weeks), "week") + Plural(int(days), "day") + Plural(int(hours), "hour") + Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else if days > 0 {
		result = Plural(int(days), "day") + Plural(int(hours), "hour") + Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else if hours > 0 {
		result = Plural(int(hours), "hour") + Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else if minutes > 0 {
		result = Plural(int(minutes), "minute") + Plural(int(seconds), "second")
	} else {
		result = Plural(int(seconds), "second")
	}

	return
}
