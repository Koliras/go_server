package api

import "net/http"

func (app App) StaticRegisterForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/register_form.html")
}

func (app App) StaticLoginForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login_form.html")
}
