package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rk1165/feedcreator/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", *dbName)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		feeds:         &models.FeedModel{DB: db},
		templateCache: templateCache,
		formDecoder:   form.NewDecoder(),
		session:       sessions.NewCookieStore([]byte("secret")),
	}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}
