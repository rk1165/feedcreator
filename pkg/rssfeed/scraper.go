package rssfeed

import (
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly"
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
	"strings"
	"time"
)

const PubDateFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func CreateFeedFile(feed *models.Feed) error {
	channel, err := CreateChannel(feed)
	if err != nil {
		return err
	}
	rss := &models.RSS{
		Version: "2.0",
		Channel: *channel,
	}
	err = internal.WriteRSSFeedFile(feed.Name, rss)
	if err != nil {
		logger.ErrorLog.Printf("error=%v", err)
		return err
	}

	logger.InfoLog.Println("RSS file created successfully")
	return nil
}

func CreateChannel(feed *models.Feed) (*models.Channel, error) {
	items := createFeedItems(feed)
	if items == nil {
		return nil, fmt.Errorf("no feed items present")
	}
	channel := &models.Channel{
		Title:       feed.Title,
		Link:        feed.Url,
		Description: feed.Description,
		Items:       items,
	}
	return channel, nil
}

func createFeedItems(feed *models.Feed) *[]models.Item {

	logger.InfoLog.Printf("creating feed items for url %s", feed.Url)
	c := colly.NewCollector(
		//colly.CacheDir("./cache"),
		//colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(30 * time.Second)
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
	if len(items) > 0 {
		return &items
	}
	logger.InfoLog.Printf("Starting dynamic website loading for creating items")
	items = createFeedItemsDynamic(feed)
	if len(items) > 0 {
		return &items
	}
	logger.WarnLog.Printf("No items were found for url=%s", feed.Url)
	return nil
}

func createFeedItemsDynamic(feed *models.Feed) []models.Item {
	logger.InfoLog.Printf("Loading URL=%s dynamically", feed.Url)
	var items []models.Item

	path := flag.Lookup("path").Value.String()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.ExecPath(path),
	)

	// Create a new context with a timeout for chromedp
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Run the chromedp tasks
	var itemNodes []*cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Navigate(feed.Url),
		chromedp.WaitVisible(feed.ItemSelector, chromedp.ByQueryAll),
		chromedp.Nodes(feed.ItemSelector, &itemNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		logger.ErrorLog.Printf("failed to fetch items for url=%s, error=%v", feed.Url, err)
		return nil
	}

	logger.InfoLog.Printf("Found %d items for url=%s", len(itemNodes), feed.Url)

	title, link, description := new(string), new(string), new(string)
	for _, itemNode := range itemNodes {

		if itemNode.NodeName == "A" {
			*link = itemNode.AttributeValue("href")
		} else {
			var linkNode []*cdp.Node
			err = chromedp.Run(ctx, chromedp.Nodes(feed.LinkSelector, &linkNode, chromedp.ByQuery, chromedp.FromNode(itemNode)))
			if err != nil {
				logger.ErrorLog.Printf("failed to fetch link for item=%s, error=%v", itemNode.NodeName, err)
				*link = ""
			} else {
				*link = linkNode[0].AttributeValue("href")
			}
		}
		err := chromedp.Run(ctx,
			chromedp.Text(feed.TitleSelector, title, chromedp.ByQuery, chromedp.FromNode(itemNode)),
		)
		if err != nil {
			logger.ErrorLog.Printf("failed to fetch title for url=%s, error=%v", feed.Url, err)
		}

		if feed.DescSelector != "" {
			err = chromedp.Run(ctx,
				chromedp.Text(feed.DescSelector, description, chromedp.ByQuery, chromedp.FromNode(itemNode)),
			)
			if err != nil {
				logger.ErrorLog.Printf("failed to fetch description for url=%s, error=%v", feed.Url, err)
			}
		}

		// Process the items
		item := models.Item{
			Title:       *title,
			Link:        *link,
			Description: *description,
			GUID:        *link,
			PubDate:     time.Now().UTC().Format(PubDateFormat),
		}
		items = append(items, item)
	}
	return items
}
