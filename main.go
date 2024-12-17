package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Koliras/go_server/api"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	s := &http.Server{
		Addr:    ":8080",
		Handler: api.Routes(db),
	}
	log.Fatal(s.ListenAndServe())
}
