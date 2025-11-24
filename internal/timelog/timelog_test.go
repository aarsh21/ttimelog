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

	// Create templates using %s where the date should go
	templates := []string{
		"%s 09:41 +0530: **arrived",
		"%s 13:30 +0530: Workin on ttimelog",
		"%s 14:15 +0530: **lunch",
		"%s 16:30 +0530: adding test for loading entries",
		"%s 17:30 +0530: working ttimelog on UI",
	}

	today := time.Now().Format("2006-01-02")

	var lines []string

	for _, tmpl := range templates {
		line := fmt.Sprintf(tmpl, today)
		lines = append(lines, line)
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
	assert.Len(t, entries, 5)

	// Check the first entry (Start of day, should be 0 duration)
	assert.Equal(t, "**arrived", strings.TrimSpace(entries[0].Description))
	assert.Equal(t, time.Duration(0), entries[0].Duration)

	// Check the second entry (Duration calculation)
	// 09:41 -> 13:30 is 3h 49m
	expectedDuration := (3 * time.Hour) + (49 * time.Minute)
	assert.Equal(t, expectedDuration, entries[1].Duration)
	assert.Contains(t, entries[1].Description, "Workin on ttimelog")
}
