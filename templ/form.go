package templ

import (
	"html/template"
	"net/http"
)

type DataRegisterForm struct {
	GeneralError         string
	Nickname             string
	NicknameErrors       []string
	Email                string
	EmailErrors          []string
	Password             string
	PasswordErrors       []string
	RepeatPassword       string
	RepeatPasswordErrors []string
}

func HtmlFormRegisterErrors(w http.ResponseWriter, data DataRegisterForm) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	template.Must(template.ParseFiles("./templ/form.html")).ExecuteTemplate(w, "register_errors", data)
}

type DataLoginForm struct {
	Error    string
	Email    string
	Password string
}

func HtmlFormLoginErrors(w http.ResponseWriter, data DataLoginForm) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	template.Must(template.ParseFiles("./templ/form.html")).ExecuteTemplate(w, "login_errors", data)
}
