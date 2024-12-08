package main

import (
	"log"
	"net/http"

	"github.com/Koliras/go_server/api"
)

func main() {
	s := &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("GET /", api.Healthcheck)
	log.Fatal(s.ListenAndServe())
}
