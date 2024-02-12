package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{fs: http.Dir("./ui/static")})
	mux.Handle("GET /static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)
	mux.HandleFunc("POST /feed/add", app.addFeed)
	mux.HandleFunc("GET /feed/view", app.viewFeed)
	mux.HandleFunc("GET /feeds", app.allFeeds)
	//mux.HandleFunc("/feed/{id}/delete", app.delete)
	//mux.HandleFunc("/feed/{id}/update", app.update)
	//mux.HandleFunc("/feed/{id}/save", app.save)

	return mux

}
