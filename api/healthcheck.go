package api

import "net/http"

func (app App) Healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am healthy"))
}
