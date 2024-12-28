package api

import (
	"net/http"

	"github.com/Koliras/go_server/templ"
)

func (app App) StaticRegisterForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/register_form.html")
}

func (app App) StaticLoginForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/login_form.html")
}

func (app App) StaticProfile(w http.ResponseWriter, r *http.Request) {
	dp := templ.DataProfile{}

	q := r.URL.Query()
	sp := q.Get("t") // get type of page content
	dp.ContentType = templ.GetContentType(sp)

	templ.HtmlProfile(w, dp)
}
