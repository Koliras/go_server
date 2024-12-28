package templ

import (
	"html/template"
	"net/http"
)

func HtmlAuthError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Add("HX-Retarget", "body")
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles("./templ/error.html")).ExecuteTemplate(w, "auth", msg)
}

// msg is optional
func HtmlInternalError(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	template.Must(template.ParseFiles("./templ/error.html")).ExecuteTemplate(w, "internal", msg)
}
