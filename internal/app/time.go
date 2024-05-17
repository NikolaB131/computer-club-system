package app

import (
	"fmt"
	"strconv"
	"strings"
)

// Parsing time string of format "XX:XX" and returns minutes
func parseTime(str string) (int, error) {
	splitted := strings.Split(str, ":")
	hours, err := strconv.Atoi(splitted[0])
	if err != nil {
		return 0, fmt.Errorf("failed to convert hours to integer: %w", err)
	}
	minutes, err := strconv.Atoi(splitted[1])
	if err != nil {
		return 0, fmt.Errorf("failed to convert minutes to integer: %w", err)
	}

	return hours*60 + minutes, nil
}

// Converts minutes to time string of format "XX:XX"
func convertMinutesToStringTime(minutes int) string {
	hours := minutes / 60
	min := minutes % 60
	hoursStr := strconv.Itoa(hours)
	minutesStr := strconv.Itoa(min)

	if hours < 10 {
		hoursStr = "0" + hoursStr
	}
	if min < 10 {
		minutesStr = "0" + minutesStr
	}

	return hoursStr + ":" + minutesStr
}
