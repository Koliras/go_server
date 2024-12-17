package api

import "net/http"

func (app App) StaticRegisterForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/register_form/index.html")
}
