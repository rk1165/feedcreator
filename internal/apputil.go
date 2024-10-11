package internal

import (
	"encoding/xml"
	"fmt"
	"github.com/rk1165/feedcreator/internal/models"
	"io"
	"os"
	"path/filepath"
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
