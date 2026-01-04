// Package timelog contains implementation required for timelogging
package timelog

import (
	"strings"
	"time"
)

type TimeLog struct {
	startTime time.Time
}

func GetTimeDiff(startTime time.Time, endTime time.Time) string {
	diff := endTime.Sub(startTime)
	return FormatDuration(diff)
}

func isSlackingTime(input string) bool {
	return strings.Contains(input, "**")
}

func parseTextInput(input string) {

}

func appendFile(input string) {

}
