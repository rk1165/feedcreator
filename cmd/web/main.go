package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
	"html/template"
	"log"
	"net/http"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	feeds         models.FeedModelInterface
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
	db            *sql.DB
	session       *sessions.CookieStore
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dbName := flag.String("db", "feeds.db", "SQLite Datasource name")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbName)
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.ErrorLog.Fatal(err)
	}

	app := &application{
		errorLog:      logger.ErrorLog,
		infoLog:       logger.InfoLog,
		feeds:         &models.FeedModel{DB: db},
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
		session:       sessions.NewCookieStore([]byte("secret")),
	}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     logger.ErrorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	internal.ScheduleFunc(60*time.Second, app.cleanFeeds)
	internal.ScheduleFunc(90*time.Second, app.updateFeeds)

	logger.InfoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	logger.ErrorLog.Fatal(err)
}
