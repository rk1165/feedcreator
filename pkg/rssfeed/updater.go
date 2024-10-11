package rssfeed

import (
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
)

// TODO : Add a channel
func UpdateFeeds(feeds []*models.Feed) {
	logger.InfoLog.Println("Started Updating feeds")
	updateFeed(feeds[0])
	logger.InfoLog.Println("Finished Updating feeds")
}

func updateFeed(feed *models.Feed) {
	logger.InfoLog.Printf("Updating feed with url %s", feed.Url)
	// Get latest items
	items := createFeedItems(feed)
	rss, err := internal.ReadRSSFeedFile(feed.Name)
	if err != nil {
		logger.ErrorLog.Printf("Error reading RSS feed: %v", err)
		return
	}

	logger.InfoLog.Printf("Current item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)
	guids := getGUIDs(*rss.Channel.Items)

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		if _, ok := guids[item.GUID]; !ok {
			*rss.Channel.Items = append([]models.Item{item}, *rss.Channel.Items...)
		}
	}
	logger.InfoLog.Printf("Updated item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)

	// Write updated XML
	err = internal.WriteRSSFeedFile(feed.Name, rss)
	if err != nil {
		logger.ErrorLog.Printf("failed to write file: err=%v", err)
		return
	}
	logger.InfoLog.Printf("Updated feed with url %s", feed.Url)
}

func getGUIDs(items []models.Item) map[string]bool {
	guids := make(map[string]bool)
	for _, item := range items {
		guids[item.GUID] = true
	}
	return guids
}
