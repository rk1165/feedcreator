package mocks

import (
	"time"

	"github.com/rk1165/feedcreator/internal/models"
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

func (m *FeedModel) Insert(feed *models.Feed) (int, error) {
	return 2, nil
}

func (m *FeedModel) GetByName(name string) (*models.Feed, error) {
	switch name {
	case "example.xml":
		return mockFeed, nil
	default:
		return nil, models.ErrNoRecord
	}
}
func (m *FeedModel) GetById(id int) (*models.Feed, error) {
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

func (m *FeedModel) Delete(id int) error {
	switch id {
	case 1:
		return nil
	default:
		return models.ErrNoRecord
	}
}
