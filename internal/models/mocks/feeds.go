package mocks

import (
	"github.com/rk1165/feedcreator/internal/models"
	"time"
)

var mockFeed = &models.Feed{
	ID:            1,
	Title:         "Example Feed",
	Url:           "https://www.example.com",
	Description:   "Feed for example",
	Name:          "example.xml",
	ItemSelector:  "div.heading",
	TitleSelector: "a.posting",
	LinkSelector:  "a.posting",
	Created:       time.Now(),
}

type FeedModel struct {
}

func (m *FeedModel) Insert(title, name, url, description, itemSelector, titleSelector,
	linkSelector, descSelector string) (int, error) {
	return 2, nil
}

func (m *FeedModel) Get(id int) (*models.Feed, error) {
	switch id {
	case 1:
		return mockFeed, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *FeedModel) All() ([]*models.Feed, error) {
	return []*models.Feed{mockFeed}, nil

}
