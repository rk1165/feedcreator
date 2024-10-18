package main

import (
	"errors"
	"fmt"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/internal/validator"
	"github.com/rk1165/feedcreator/pkg/logger"
	"github.com/rk1165/feedcreator/pkg/rssfeed"
	"net/http"
	"os"
	"strconv"
	"time"
)

type feedCreateForm struct {
	Title               string `form:"title"`
	URL                 string `form:"url"`
	Description         string `form:"desc"`
	Name                string `form:"name"`
	ItemSelector        string `form:"item_elem"`
	TitleSelector       string `form:"title_elem"`
	LinkSelector        string `form:"link_elem"`
	DescSelector        string `form:"desc_elem"`
	validator.Validator `form:"-"`
}

// home for rendering home page
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	app.render(w, http.StatusOK, "home.tmpl", data)
}

// createFeed for displaying the feed creation form
func (app *application) createFeed(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	data.Form = feedCreateForm{}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

// submitFeed for submitting the feed
func (app *application) submitFeed(w http.ResponseWriter, r *http.Request) {
	var form feedCreateForm

	err := app.decodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters")
	form.CheckField(validator.NotBlank(form.URL), "url", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.ItemSelector), "itemSelector", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.TitleSelector), "titleSelector", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.LinkSelector), "linkSelector", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(w, r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	feed, err := app.feeds.GetByName(form.Name + ".xml")
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.infoLog.Printf("Feed %s does not exist", form.Name)
		} else {
			app.serverError(w, err)
			return
		}
	}

	if feed != nil {
		app.infoLog.Printf("feed %s already exists. Redirecting to it", form.Name)
		http.Redirect(w, r, fmt.Sprintf("/feed/view?id=%d", feed.ID), http.StatusSeeOther)
		return
	}

	feed = &models.Feed{
		Title:         form.Title,
		Name:          form.Name + ".xml",
		Url:           form.URL,
		Description:   form.Description,
		ItemSelector:  form.ItemSelector,
		TitleSelector: form.TitleSelector,
		LinkSelector:  form.LinkSelector,
		DescSelector:  form.DescSelector,
	}

	session, _ := app.session.Get(r, "flash-session")

	go func() {
		err := rssfeed.CreateFeedFile(feed)
		if err != nil {
			logger.ErrorLog.Printf("Failed to create RSS feed for url=%s", feed.Url)
			app.serverError(w, err)
			return
		}
		id, err := app.feeds.Insert(feed)
		if err != nil {
			app.serverError(w, err)
			return
		}
		app.infoLog.Printf("Successfully created RSS feed with id=%d for url=%s", id, feed.Url)
	}()

	session.Values["flash"] = fmt.Sprintf("Creating feed for %s", form.URL)
	session.Save(r, w)

	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

// viewFeed for viewing a specific feed
func (app *application) viewFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	feed, err := app.feeds.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(w, r)
	data.Feed = feed
	app.render(w, http.StatusOK, "view.tmpl", data)
}

// deleteFeed for deleting a feed
func (app *application) deleteFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	feed, _ := app.feeds.GetById(id)
	err = app.feeds.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = os.Remove(fmt.Sprintf("./ui/static/rss/%s", feed.Name))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Printf("removed feed file %s", feed.Name)

	session, _ := app.session.Get(r, "flash-session")
	session.Values["flash"] = fmt.Sprintf("Deleted feed with id=%s", feed.Name)
	_ = session.Save(r, w)

	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

// allFeeds gets list of all the feeds
func (app *application) allFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.feeds.All()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(w, r)
	data.Feeds = feeds

	app.render(w, http.StatusOK, "feeds.tmpl", data)
}

// updateFeeds updates every feed at a scheduler interval
func (app *application) updateFeeds(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("Running updateFeeds at %v", time.Now().UTC())
	feeds, err := app.feeds.All()
	if err != nil {
		app.serverError(w, err)
		return
	}
	rssfeed.UpdateFeeds(feeds)
}

// cleanFeeds cleans feed entries which are more than 3 days old
func (app *application) cleanFeeds(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Printf("Running cleanFeeds at %v", time.Now().UTC())
	feeds, err := app.feeds.All()
	if err != nil {
		app.serverError(w, err)
		return
	}
	rssfeed.CleanFeeds(feeds)
}
