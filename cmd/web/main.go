package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/feedcreator/internal"
	"github.com/rk1165/feedcreator/internal/models"
	"github.com/rk1165/feedcreator/pkg/logger"
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
	_ = flag.String("path", "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser",
		"Path to browser for executing dynamic content")
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

	// Schedules cleaning and updating feeds
	internal.ScheduleFunc(90*time.Minute, app.cleanFeeds)
	internal.ScheduleFunc(60*time.Minute, app.updateFeeds)

	logger.InfoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	logger.ErrorLog.Fatal(err)
}
