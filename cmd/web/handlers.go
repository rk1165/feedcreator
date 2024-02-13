package main

import (
	"errors"
	"fmt"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/internal/scraper"
	"github.com/rk1165/feedcreator/internal/validator"
	"net/http"
	"strconv"
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) feedCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(w, r)
	data.Form = feedCreateForm{}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) feedCreatePost(w http.ResponseWriter, r *http.Request) {
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

	feed, err := app.feeds.GetByName(form.Name)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.infoLog.Printf("Feed %s does not exist", form.Name)
		} else {
			app.serverError(w, err)
		}
	}

	if feed != nil {
		app.infoLog.Printf("feed with name %s already exists", form.Name)
		http.Redirect(w, r, "/feeds/", http.StatusSeeOther)
	}

	feed = &models.Feed{Title: form.Title, Name: form.Name + ".xml", Url: form.URL, Description: form.Description,
		ItemSelector: form.ItemSelector, TitleSelector: form.TitleSelector, LinkSelector: form.LinkSelector,
		DescSelector: form.DescSelector}

	scraper.CreateFeedFile(feed)

	id, err := app.feeds.Insert(feed)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Printf("created feed with id %d", id)
	// Redirect the user to the relevant page for the feed
	http.Redirect(w, r, fmt.Sprintf("/feed/view/%d", id), http.StatusSeeOther)
}

func (app *application) viewFeed(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("GET params were:", r.URL.Path)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
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

// TODO : Create go workers for fetching and updating the feeds at regular interval
