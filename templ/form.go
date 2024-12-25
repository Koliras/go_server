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
	PasswordErrors       []string
	Password             string
	RepeatPasswordErrors []string
	RepeatPassword       string
}

func FormRegisterErrors(w http.ResponseWriter, data RegisterFormData) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	template.Must(template.ParseFiles("./templ/form.html")).ExecuteTemplate(w, "register_errors", data)
}
