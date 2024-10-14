package main

import (
	"github.com/justinas/alice"
	"github.com/rk1165/feedcreator/ui"
	"io/fs"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	staticFiles, _ := fs.Sub(ui.Files, "static")
	fileServer := http.FileServer(http.FS(staticFiles))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	feedServer := http.FileServer(http.Dir("./rss"))
	mux.Handle("GET /rss/", http.StripPrefix("/rss/", feedServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("GET /feed/create", app.feedCreate)
	mux.HandleFunc("POST /feed/create", app.feedCreatePost)
	mux.HandleFunc("GET /feed/view/", app.viewFeed)
	mux.HandleFunc("GET /feed/delete/", app.deleteFeed)
	mux.HandleFunc("GET /feeds", app.allFeeds)
	mux.HandleFunc("GET /update", app.updateFeeds)
	mux.HandleFunc("GET /clean", app.cleanFeeds)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
