package main

import (
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

// templateData holds any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Feed        *models.Feed
	Feeds       []*models.Feed
	Form        any
	Flash       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}
		// Parse the base template file into a template set.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
