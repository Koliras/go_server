package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Koliras/go_server/middleware"
	"github.com/Koliras/go_server/templ"
	"github.com/Koliras/go_server/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (app App) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	registerFormData := templ.DataRegisterForm{}

	err := r.ParseForm()
	if err != nil {
		registerFormData.GeneralError = "Failed to parse body of request"
		templ.HtmlFormRegisterErrors(w, registerFormData)
		return
	}

	data := utils.RegisterInput{
		Nickname:       r.PostFormValue("nickname"),
		Email:          r.PostFormValue("email"),
		Password:       r.PostFormValue("password"),
		RepeatPassword: r.PostFormValue("repeat_password"),
	}

	registerFormData = utils.ValidateRegisterData(data)
	registerFormData.Nickname = data.Nickname
	registerFormData.Email = data.Email
	registerFormData.Password = data.Password
	registerFormData.RepeatPassword = data.Password

	var checkError error = nil
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		used, err := app.DB.UserExists(&data.Email)
		if err != nil {
			checkError = err
			wg.Done()
			return
		}
		if used {
			registerFormData.EmailErrors = append(registerFormData.EmailErrors, "This email is already in use")
		}
		wg.Done()
	}()

	used, err := app.DB.NicknameUsed(&data.Nickname)
	if err != nil {
		checkError = err
	}
	if used {
		registerFormData.NicknameErrors = append(registerFormData.NicknameErrors, "This nickname is already in use")
	}

	wg.Wait()

	if checkError != nil {
		registerFormData.GeneralError = "Internal server error"
		templ.HtmlFormRegisterErrors(w, registerFormData)
		return
	}

	hashedPassword, err := utils.HashString(&data.Password)
	if err != nil {
		registerFormData.GeneralError = "Internal server error"
		templ.HtmlFormRegisterErrors(w, registerFormData)
		return
	}

	if len(registerFormData.GeneralError) != 0 ||
		len(registerFormData.PasswordErrors) != 0 ||
		len(registerFormData.RepeatPasswordErrors) != 0 ||
		len(registerFormData.EmailErrors) != 0 ||
		len(registerFormData.NicknameErrors) != 0 {

		templ.HtmlFormRegisterErrors(w, registerFormData)
		return
	}

	err = app.DB.CreateUser(&data.Nickname, &data.Email, &hashedPassword)
	if err != nil {
		registerFormData.GeneralError = "Error when creating user"
		templ.HtmlFormRegisterErrors(w, registerFormData)
		return
	}
	w.Header().Add("HX-Push-Url", "/profile")
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app App) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	loginFormData := templ.DataLoginForm{}

	err := r.ParseForm()
	if err != nil {
		loginFormData.Error = "Failed to parse the body of request"
		templ.HtmlFormLoginErrors(w, loginFormData)
		return
	}

	data := loginBody{
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	loginFormData.Email = data.Email
	loginFormData.Password = data.Password

	user, err := app.DB.GetUserByEmail(&data.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			loginFormData.Error = "User with such email not found"
			templ.HtmlFormLoginErrors(w, loginFormData)
			return
		}
		loginFormData.Error = "Internal server error"
		templ.HtmlFormLoginErrors(w, loginFormData)
		return
	}

	isSamePass := utils.CompareStrWithHash(&data.Password, &user.Password)
	if !isSamePass {
		loginFormData.Error = "Incorrect email or password"
		templ.HtmlFormLoginErrors(w, loginFormData)
		return
	}

	claims := middleware.JwtClaims{
		user.Email,
		user.Nickname,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(middleware.JwtKey)
	if err != nil {
		loginFormData.Error = "Incorrect email or password"
		templ.HtmlFormLoginErrors(w, loginFormData)
		return
	}
	if len(loginFormData.Error) == 0 {
		templ.HtmlFormLoginErrors(w, loginFormData)
		return
	}
	cookie := &http.Cookie{
		Name:     middleware.JwtTokenCookieName,
		Value:    token,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
	w.Header().Add("HX-Push-Url", "/profile")
}

func (app App) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.AllUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	strUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(strUsers)
}
