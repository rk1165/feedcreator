package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rk1165/feedcreator/pkg/logger"
)

type FeedModelInterface interface {
	Insert(feed *Feed) (int, error)
	GetByName(name string) (*Feed, error)
	GetById(id int) (*Feed, error)
	Delete(id int) error
	All() ([]*Feed, error)
}

// Feed type to hold the data for an individual feed.
type Feed struct {
	ID            int
	Title         string
	Url           string
	Description   string
	Name          string
	ItemSelector  string
	TitleSelector string
	LinkSelector  string
	DescSelector  string
	Created       time.Time
}

// FeedModel a type which wraps a sql.DB connection pool
type FeedModel struct {
	DB *sql.DB
}

// Insert will insert a Feed into the DB
func (m *FeedModel) Insert(feed *Feed) (int, error) {
	stmt := `INSERT INTO feed (title, name, url, description,
			 item_selector, title_selector,  link_selector, desc_selector)
			 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, feed.Title, feed.Name, feed.Url, feed.Description,
		feed.ItemSelector, feed.TitleSelector, feed.LinkSelector, feed.DescSelector)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *FeedModel) GetByName(name string) (*Feed, error) {
	stmt := `SELECT id, title, name, url, description, item_selector, title_selector,
			 link_selector, desc_selector, created
			 FROM feed WHERE name = ?`

	row := m.DB.QueryRow(stmt, name)

	feed := &Feed{}
	err := row.Scan(&feed.ID, &feed.Title, &feed.Name, &feed.Url, &feed.Description,
		&feed.ItemSelector, &feed.TitleSelector,
		&feed.LinkSelector, &feed.DescSelector, &feed.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return feed, nil
}

func (m *FeedModel) GetById(id int) (*Feed, error) {
	stmt := `SELECT id, title, name, url, description, item_selector, title_selector,
			 link_selector, desc_selector, created
			 FROM feed WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	feed := &Feed{}
	err := row.Scan(&feed.ID, &feed.Title, &feed.Name, &feed.Url, &feed.Description,
		&feed.ItemSelector, &feed.TitleSelector,
		&feed.LinkSelector, &feed.DescSelector, &feed.Created)

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
			 item_selector, title_selector, link_selector, desc_selector, created
			 FROM feed ORDER BY id DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*Feed

	for rows.Next() {
		feed := &Feed{}
		err = rows.Scan(&feed.ID, &feed.Title, &feed.Name, &feed.Url, &feed.Description,
			&feed.ItemSelector, &feed.TitleSelector,
			&feed.LinkSelector, &feed.DescSelector, &feed.Created)

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

func (m *FeedModel) Delete(id int) error {
	stmt := `DELETE FROM feed WHERE id = ?`
	rows, err := m.DB.Exec(stmt, id)
	if err != nil {
		err := fmt.Errorf("failed to delete feed with id %d: %v", id, err)
		return err
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		err := fmt.Errorf("failed to delete feed with id %d: %v", id, err)
		return err
	}
	logger.InfoLog.Printf("%d rows deleted", affected)
	return nil
}
