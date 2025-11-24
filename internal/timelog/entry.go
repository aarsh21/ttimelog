package timelog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Entry struct {
	EndTime     time.Time
	Description string
	// Duration is computed on load, not stored
	Duration time.Duration
}

const timeLayout = "2006-01-02 15:04 -0700"

// SaveEntry saves the entry in 'YYYY-MM-DD HH:MM +/-0000: Task Description' format
func SaveEntry(entry Entry) error {
	// TODO: make it dynamic
	// we already create the file if didn't exist in config.go when we create
	// we can save the path of the file in config struct or make a global variable.
	// I am not sure.
	filename := "/home/rashesh/.ttimelog/ttimelog.txt"

	// Open the file in append mode. Create it if it doesn't exist.
	// os.O_APPEND: Open the file for appending.
	// os.O_WRONLY: Open the file for writing only.
	// 0644: File permissions (read/write for owner, read-only for others).
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	dateAndTime := entry.EndTime.Format(timeLayout)
	textEntry := fmt.Sprintf("%s: %s\n", dateAndTime, entry.Description)

	if _, err := f.WriteString(textEntry); err != nil {
		return err
	}

	return nil
}

// 2025-10-17 13:30 +0530: Workin on ttimelog
func parseEntry(line string, previousEntryEndTime *time.Time) (*Entry, error) {
	// It splits in 3 strings and we merge them later
	tokens := strings.Split(line, ":")
	if len(tokens) < 3 {
		return nil, errors.New("invalid format")
	}

	dateAndTime := tokens[0] + ":" + tokens[1]
	dateAndTimeTokens := strings.Split(dateAndTime, " ")
	if len(dateAndTimeTokens) < 3 {
		return nil, errors.New("invalid format")
	}

	parsedDate := dateAndTimeTokens[0]
	currentDate := time.Now().Format("2006-01-02")

	if parsedDate != currentDate {
		return nil, nil
	}

	endTime, err := time.Parse(timeLayout, dateAndTime)
	if err != nil {
		return nil, err
	}
	entryDuration, _ := time.ParseDuration("0s")
	if previousEntryEndTime != nil {
		entryDuration = endTime.Sub(*previousEntryEndTime)
	}
	return &Entry{
		EndTime:     endTime,
		Description: tokens[2],
		Duration:    entryDuration,
	}, nil
}

func LoadEntries(filePath string) ([]Entry, error) {
	entries := make([]Entry, 0)
	file, err := os.Open(filePath)
	if err != nil {
		return entries, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var ( entry *Entry
			err   error
		)
		line = strings.Trim(line, " ")
		if len(entries) == 0 {
			entry, err = parseEntry(line, nil)
		} else {
			entry, err = parseEntry(line, &entries[len(entries)-1].EndTime)
		}
		if err != nil {
			return entries, err
		}

		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return entries, err
	}
	return entries, nil
}
