package scraper

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/gocolly/colly"
	"github.com/rk1165/feedcreator/internal/models"
	"log"
	"strings"
	"time"
)

func createFeedItems(feed *models.Feed) []Item {

	log.Printf("creating feed items for url %s", feed.Url)
	c := colly.NewCollector(
		//colly.CacheDir("./cache"),
		//colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)

	var items []Item

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
		hashedLink := sha256.Sum256([]byte(link))
		guid := base64.URLEncoding.EncodeToString(hashedLink[:])
		desc := e.DOM.Find(feed.DescSelector).Text()
		item := Item{
			Title:       title,
			Link:        link,
			Description: desc,
			GUID:        guid,
			IsPermalink: false,
			PubDate:     time.Now().UTC().Format("2006-01-02 15:04:05.000000"),
		}
		items = append(items, item)
	})

	err := c.Visit(feed.Url)
	if err != nil {
		log.Println(err)
		return nil
	}
	return items
}

func CreateChannel(feed *models.Feed) *Channel {
	items := createFeedItems(feed)
	channel := &Channel{Title: feed.Title, Link: feed.Url, Description: feed.Description,
		Items: &items}
	return channel
}
