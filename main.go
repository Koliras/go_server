package main

import (
	"fmt"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", RequestHandler)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println("Failed to initialize server. Error:", err)
	}
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request from %s\n", r.URL.Path)
}
