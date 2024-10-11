package rssfeed

import (
	"github.com/gocolly/colly"
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
	"strings"
	"time"
)

const PubDateFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func CreateFeedFile(feed *models.Feed) {
	channel := CreateChannel(feed)
	rss := &models.RSS{
		Version: "2.0",
		Channel: *channel,
	}
	err := internal.WriteRSSFeedFile(feed.Name, rss)
	if err != nil {
		logger.ErrorLog.Printf("error=%v", err)
		return
	}

	logger.InfoLog.Println("RSS file created successfully")
}

func CreateChannel(feed *models.Feed) *models.Channel {
	items := createFeedItems(feed)
	channel := &models.Channel{
		Title:       feed.Title,
		Link:        feed.Url,
		Description: feed.Description,
		Items:       &items,
	}
	return channel
}

func createFeedItems(feed *models.Feed) []models.Item {

	logger.InfoLog.Printf("creating feed items for url %s", feed.Url)
	c := colly.NewCollector(
		//colly.CacheDir("./cache"),
		//colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)

	var items []models.Item

	c.OnHTML(feed.ItemSelector, func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.DOM.Find(feed.TitleSelector).First().Text())
		var link string
		if e.Name == "a" {
			link = strings.TrimSpace(e.Attr("href"))
		} else {
			link = strings.TrimSpace(e.DOM.Find(feed.LinkSelector).AttrOr("href", ""))
		}
		if link != "" {
			link = e.Request.AbsoluteURL(link)
		}
		desc := e.DOM.Find(feed.DescSelector).Text()
		item := models.Item{
			Title:       title,
			Link:        link,
			Description: desc,
			GUID:        link,
			PubDate:     time.Now().UTC().Format(PubDateFormat),
		}
		items = append(items, item)
	})

	err := c.Visit(feed.Url)
	if err != nil {
		logger.ErrorLog.Printf("failed to fetch feed items for url=%s, error=%v", feed.Url, err)
		return nil
	}
	return items
}
