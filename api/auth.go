package api

import (
	"database/sql"
	"encoding/json"
	"io"
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
	registerFormData := templ.RegisterFormData{}

	err := r.ParseForm()
	if err != nil {
		registerFormData.GeneralError = "Failed to parse body of request"
		templ.FormRegisterErrors(w, registerFormData)
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
		templ.FormRegisterErrors(w, registerFormData)
		return
	}

	hashedPassword, err := utils.HashString(&data.Password)
	if err != nil {
		registerFormData.GeneralError = "Internal server error"
		templ.FormRegisterErrors(w, registerFormData)
		return
	}

	if len(registerFormData.GeneralError) != 0 ||
		len(registerFormData.PasswordErrors) != 0 ||
		len(registerFormData.RepeatPasswordErrors) != 0 ||
		len(registerFormData.EmailErrors) != 0 ||
		len(registerFormData.NicknameErrors) != 0 {

		templ.FormRegisterErrors(w, registerFormData)
		return
	}

	err = app.DB.CreateUser(&data.Nickname, &data.Email, &hashedPassword)
	if err != nil {
		registerFormData.GeneralError = "Error when creating user"
		templ.FormRegisterErrors(w, registerFormData)
		return
	}
	w.Header().Add("HX-Push-Url", "/profile")
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app App) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error when reading the body of the request"))
		return
	}

	data := loginBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user, err := app.DB.GetUserByEmail(&data.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User with such email not found", 404)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}

	isSamePass := utils.CompareStrWithHash(&data.Password, &user.Password)
	if !isSamePass {
		http.Error(w, "Incorrect email or password", http.StatusNotFound)
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
		http.Error(w, err.Error(), 500)
		return
	}
	cookie := &http.Cookie{
		Name:     middleware.JwtTokenCookieName,
		Value:    token,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
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
