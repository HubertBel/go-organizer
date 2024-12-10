package utils

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
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

func ValidateTime(value string) bool {
	regex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$`)
	
	if !regex.MatchString(value) {
		return false
	}
	
	parts := strings.Split(value, " ")
	dateParts := strings.Split(parts[0], "-")
	timeParts := strings.Split(parts[1], ":")
	
	year, err := strconv.Atoi(dateParts[0])
	if err != nil || year <= 0 {
		return false
	}
	
	month, err := strconv.Atoi(dateParts[1])
	if err != nil || month <= 0 || month > 12 {
		return false
	}
	
	day, err := strconv.Atoi(dateParts[2])
	if err != nil || day <= 1 || day > 31 {
		return false
	}
	
	hours, err := strconv.Atoi(timeParts[0])
	if err != nil || hours < 0 || hours > 23 {
		return false
	}
	
	minutes, err := strconv.Atoi(timeParts[1])
	if err != nil || (minutes != 0 && minutes != 30) {
		return false
	}
	
	_, err = time.Parse("2006-01-02 15:04", value)
	return err == nil
}

func ValidateName(value string) bool {
    if value == "" {
        return false
    }

    return true
}

func ValidateNumber(value string) bool {
    n, err := strconv.Atoi(value)
    if err != nil {
        return false
    }

    if n <= 0 {
        return false
    }

    return true
}

func ValidateDuration(value string) bool {
    duration, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return false
    }

    if duration <= 0.0 {
        return false
    }

    if math.Mod(duration, 0.5) != 0 {
         return false
    }

    return true
}
