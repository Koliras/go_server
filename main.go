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
	http.HandleFunc("POST /auth/register", api.Register)
	log.Fatal(s.ListenAndServe())
}
