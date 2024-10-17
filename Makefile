run:
	go run ./cmd/web/

init:
	sqlite3 feeds.db < ./sql/ddl.sql
	go mod tidy

build:
	go build -o feedcreator ./cmd/web/

clean:
	sqlite3 feeds.db < ./sql/clean.sql

.PHONY: run init build clean