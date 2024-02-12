package main

import (
	"errors"
	"fmt"
	"github.com/rk1165/feedcreator/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	templateSet, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = templateSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) addFeed(w http.ResponseWriter, r *http.Request) {
	title := "No Fluff Jobs"
	name := "nofluffjobs.xml"
	url := "https://nofluffjobs.com/"
	description := "No Fluff Jobs Feed"
	itemTag := "a"
	itemCls := "posting-list-item"
	titleTag := "h3"
	titleCls := "posting-title__position"
	linkTag := "a"
	linkCls := "posting-list-item"
	descriptionTag := ""
	descriptionCls := ""

	id, err := app.feeds.Insert(title, name, url, description, itemTag, itemCls,
		titleTag, titleCls, linkTag, linkCls, descriptionTag, descriptionCls)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Redirect the user to the relevant page for the feed
	http.Redirect(w, r, fmt.Sprintf("/feed/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) viewFeed(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	feed, err := app.feeds.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Fprintf(w, "%+v", feed)
}

func (app *application) allFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.feeds.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, feed := range feeds {
		fmt.Fprintf(w, "%+v\n", feed)
	}

}
