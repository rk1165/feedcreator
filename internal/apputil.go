package internal

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rk1165/feedcreator/internal/models"
)

func ReadRSSFeedFile(feedName string) (*models.RSS, error) {

	filePath := filepath.Join("./ui/static/rss", feedName)
	file, err := os.Open(filePath)
	if err != nil {
		err := fmt.Errorf("failed to open rss feed file: %v", err)
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		err := fmt.Errorf("failed to read rss feed file: %v", err)
		return nil, err
	}

	rss := new(models.RSS)
	err = xml.Unmarshal(bytes, rss)
	if err != nil {
		err := fmt.Errorf("failed to unmarshal rss feed file: %v", err)
		return nil, err
	}

	return rss, nil
}

func WriteRSSFeedFile(feedName string, rss *models.RSS) error {
	output, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		err := fmt.Errorf("failed to marshal rss feed file: %v", err)
		return err
	}
	filePath := filepath.Join("./ui/static/rss", feedName)
	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		err := fmt.Errorf("failed to write rss feed file: %v", err)
		return err
	}
	return nil
}

func ScheduleFunc(interval time.Duration, task func(w http.ResponseWriter, r *http.Request)) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				task(nil, nil)
			}
		}
	}()
}
