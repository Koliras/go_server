package templ

import (
	"html/template"
	"net/http"
)

func GetContentType(s string) string {
	switch s {
	case "p":
		return "posts"
	case "s":
		return "settings"
	default:
		return "posts"
	}
}

type DataProfile struct {
	Nickname    string
	ContentType string
	Posts       any
	Settings    any
}

func HtmlProfile(w http.ResponseWriter, p DataProfile) {
	template.Must(template.ParseFiles("./templ/profile.html")).ExecuteTemplate(w, "main", p)
}

type DataSettings struct {
}

func HtmlProfileSettings(w http.ResponseWriter, settings DataSettings) {
	template.Must(template.ParseFiles("./templ/profile.html")).ExecuteTemplate(w, "settings", settings)
}

type DataPosts struct {
	Posts []any
}

func HtmlProfilePosts(w http.ResponseWriter, p DataPosts) {
	template.Must(template.ParseFiles("./templ/profile.html")).ExecuteTemplate(w, "posts", p)
}
