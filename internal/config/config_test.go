package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupTimeLogDirectory(t *testing.T) {
	tempDir := t.TempDir()

	_, err := SetupTimeLogDirectory(tempDir)
	if err != nil {
		t.Fatalf("SetupTimeLogDirectory() failed: %v", err)
	}

	expectedFilePath := filepath.Join(tempDir, TimeLogDirname, TimeLogFilename)
	if _, err := os.Stat(expectedFilePath); err != nil {
		t.Errorf("Expected file to be created, but got error: %v", err)
	}
}
