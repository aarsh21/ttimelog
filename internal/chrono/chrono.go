package chrono

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Rash419/ttimelog/internal/config"
)

func FetchProjectList(appConfig *config.AppConfig) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", appConfig.Gtimelog.TaskListURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", appConfig.Gtimelog.AuthHeader)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	slog.Debug("Response", "status", resp.Status, "body", string(body))

	projectListPath := filepath.Join(appConfig.TimeLogDirPath, config.ProjectListFile)
	projectListFile, err := os.Create(projectListPath)
	if err != nil {
		return fmt.Errorf("failed to create project-list[%s] with error[%v]", projectListPath, err)
	}
	_, err = projectListFile.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write to project-list[%s] with error[%v]", projectListPath, err)
	}

	return nil
}
