package utils

import (
	"fmt"
	"strings"
	"time"
)

func DurationToHeight(d float64) int {
	return int(d * 2)
}

func FormatDate(t time.Time) string {
    return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func FormatHourFromTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
}

func FormatHour(hour, minute int) string {
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

func TimeToPosition(t time.Time, s string) int {

    time := FormatHourFromTime(t)
	lines := strings.Split(s, "\n")

	for i, v := range lines {
		if strings.Contains(v, time) {
			return i
		}
	}

	return -1
}
