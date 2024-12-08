package api

import "net/http"

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am healthy"))
}
