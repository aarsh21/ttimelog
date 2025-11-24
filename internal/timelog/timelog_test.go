package timelog

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeLog(t *testing.T) {
	startTime := time.Now()
	endTime := startTime.Add(2*time.Hour + 1*time.Minute)
	timeDiff := GetTimeDiff(startTime, endTime)
	assert.Equal(t, "2 h 1 min", timeDiff)
}

func TestLoadEntries(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "ttimelog.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	today := time.Now().Format("2006-01-02")

	lines := []string{
		fmt.Sprintf("%s 22:00 +0530: Yesterday task", yesterday),
		// Yesterday's last task
		fmt.Sprintf("%s 23:00 +0530: End of yesterday", yesterday),
		// Today's first task (Gap should be ignored)
		fmt.Sprintf("%s 09:00 +0530: Start of today", today),
		fmt.Sprintf("%s 10:00 +0530: Working", today),
	}

	result := strings.Join(lines, "\n")
	result = strings.TrimRight(result, "\n")

	if _, err := tmpFile.WriteString(result); err != nil {
		t.Fatalf("Failed write content to temp file with error[%v]", err)
	}

	tmpFilename := tmpFile.Name()
	tmpFile.Close()

	entries, err := LoadEntries(tmpFilename)

	assert.NoError(t, err)
	assert.Len(t, entries, 4)

	// Assertions
	// Entry 0 -> Duration 0 (first entry in timelog)
	assert.Equal(t, time.Duration(0), entries[0].Duration)
	// Entry 1 -> Duration 1 h (Yesterday's last task)
	assert.Equal(t, 1*time.Hour, entries[1].Duration)
	// Entry 2 (Today 09:00) -> Duration 0 (Reset! Not 10 hours)
	assert.Equal(t, time.Duration(0), entries[2].Duration)
	// Entry 3 (Today 10:00) -> Duration 1h
	assert.Equal(t, 1*time.Hour, entries[3].Duration)
}
