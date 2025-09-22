package rssfeed

import (
	"sync"

	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
)

func UpdateFeeds(feeds []*models.Feed) {
	logger.InfoLog.Println("Started Updating all feeds")
	var wg sync.WaitGroup
	for _, feed := range feeds {
		wg.Add(1)
		go updateFeed(feed, &wg)
	}
	wg.Wait()
	logger.InfoLog.Println("Finished Updating all feeds")
}

func updateFeed(feed *models.Feed, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.InfoLog.Printf("Updating feed with url %s", feed.Url)
	// Get latest items
	items := createFeedItems(feed)
	rss, err := internal.ReadRSSFeedFile(feed.Name)
	if err != nil {
		logger.ErrorLog.Printf("Error reading RSS feed: %v", err)
		return
	}

	logger.InfoLog.Printf("Before updating item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)
	guids := getGUIDs(*rss.Channel.Items)

	for _, item := range *items {
		if _, ok := guids[item.GUID]; !ok {
			*rss.Channel.Items = append(*rss.Channel.Items, item)
		}
	}
	logger.InfoLog.Printf("After updating item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)

	// Write updated XML
	err = internal.WriteRSSFeedFile(feed.Name, rss)
	if err != nil {
		logger.ErrorLog.Printf("failed to write file: err=%v", err)
		return
	}
	logger.InfoLog.Printf("Updated feed=%s with url %s", feed.Name, feed.Url)
}

func getGUIDs(items []models.Item) map[string]bool {
	guids := make(map[string]bool)
	for _, item := range items {
		guids[item.GUID] = true
	}
	return guids
}
