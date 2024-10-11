package rssfeed

import (
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
	"time"
)

func CleanFeeds(feeds []*models.Feed) {
	logger.InfoLog.Println("Started Cleaning feeds")
	cleanFeed(feeds[0])
	logger.InfoLog.Println("Finished Cleaning feeds")
}

func cleanFeed(feed *models.Feed) {

	logger.InfoLog.Printf("Cleaning feed with url=%s", feed.Url)
	rss, err := internal.ReadRSSFeedFile(feed.Name)
	if err != nil {
		logger.ErrorLog.Printf("Error reading RSS feed file, err=%v", err)
		return
	}
	logger.InfoLog.Printf("Current item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)
	var newItems []models.Item

	for _, item := range *rss.Channel.Items {
		if !isItemDaysOld(item, 3) {
			newItems = append(newItems, item)
		}
	}
	rss.Channel.Items = &newItems
	logger.InfoLog.Printf("Updated item_counts=%d for feed=%s", len(*rss.Channel.Items), feed.Name)
	err = internal.WriteRSSFeedFile(feed.Name, rss)
	if err != nil {
		logger.ErrorLog.Printf("Error writing RSS feed file, err=%v", err)
		return
	}
	logger.InfoLog.Printf("Cleaned feed with url=%s", feed.Url)

}

func isItemDaysOld(item models.Item, days int) bool {
	t, err := time.Parse(PubDateFormat, item.PubDate)
	if err != nil {
		logger.ErrorLog.Printf("Error parsing date: %v", err)
		return false
	}
	return time.Since(t) > (time.Duration(days) * 24 * time.Hour)
}
