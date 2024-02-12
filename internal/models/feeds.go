package models

import (
	"database/sql"
	"errors"
	"time"
)

// Feed type to hold the data for an individual feed.
type Feed struct {
	ID             int
	Title          string
	Name           string
	URL            string
	Description    string
	ItemTag        string
	ItemCls        string
	TitleTag       string
	TitleCls       string
	LinkTag        string
	LinkCls        string
	DescriptionTag string
	DescriptionCls string
	Created        time.Time
}

// FeedModel a type which wraps a sql.DB connection pool
type FeedModel struct {
	DB *sql.DB
}

// Insert will insert a Feed into the DB
func (m *FeedModel) Insert(title, name, url, description, itemTag, itemCls,
	titleTag, titleCls, linkTag, linkCls, descriptionTag, descriptionCls string) (int, error) {
	stmt := `INSERT INTO feed (title, name, url, description,
			 item_tag, item_cls, title_tag, title_cls, link_tag, link_cls, description_tag, description_cls, created)
			 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, name, url, description, itemTag, itemCls,
		titleTag, titleCls, linkTag, linkCls, descriptionTag, descriptionCls)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *FeedModel) Get(id int) (*Feed, error) {
	stmt := `SELECT id, title, name, url, description,
			 item_tag, item_cls, title_tag, title_cls, link_tag, link_cls, description_tag, description_cls
			 FROM feed WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	feed := &Feed{}
	err := row.Scan(&feed.ID, &feed.Title, &feed.Name, &feed.URL, &feed.Description,
		&feed.ItemTag, &feed.ItemCls, &feed.TitleTag, &feed.TitleCls,
		&feed.LinkTag, &feed.LinkCls, &feed.DescriptionTag, &feed.DescriptionCls)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return feed, nil
}

func (m *FeedModel) All() ([]*Feed, error) {
	stmt := `SELECT id, title, name, url, description,
			 item_tag, item_cls, title_tag, title_cls, link_tag, link_cls, description_tag, description_cls
			 FROM feed ORDER BY id DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feeds := []*Feed{}

	for rows.Next() {
		feed := &Feed{}
		err = rows.Scan(&feed.ID, &feed.Title, &feed.Name, &feed.URL, &feed.Description,
			&feed.ItemTag, &feed.ItemCls, &feed.TitleTag, &feed.TitleCls,
			&feed.LinkTag, &feed.LinkCls, &feed.DescriptionTag, &feed.DescriptionCls)

		if err != nil {
			return nil, err
		}
		feeds = append(feeds, feed)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return feeds, nil
}

func (m *FeedModel) Update() (*Feed, error) {
	return nil, nil
}

func (m *FeedModel) Delete() error {
	return nil
}
