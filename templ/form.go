package templ

import (
	"html/template"
	"net/http"
)

type RegisterFormData struct {
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

func FormRegisterErrors(w http.ResponseWriter, data RegisterFormData) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	template.Must(template.ParseFiles("./templ/form.html")).ExecuteTemplate(w, "register_errors", data)
}

type LoginFormData struct {
	Error    string
	Email    string
	Password string
}

func FormLoginErrors(w http.ResponseWriter, data LoginFormData) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	template.Must(template.ParseFiles("./templ/form.html")).ExecuteTemplate(w, "login_errors", data)
}
