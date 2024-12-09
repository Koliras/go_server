package api

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"unicode"
)

func IsValidPassword(p *string) (bool, string) {
	if len(*p) < 8 {
		return false, "Password is too short"
	}
	type condition struct {
		valid   bool
		message string
	}
	var (
		hasNum       = condition{false, "Password has to contain at least 1 number"}
		hasUpperChar = condition{false, "Password has to contain at least 1 uppercase letter"}
		hasLowerChar = condition{false, "Password has to contain at least 1 lowercase letter"}
		hasSymbol    = condition{false, "Password has to contain at least 1 symbol"}
	)
	for _, char := range *p {
		switch {
		case unicode.IsDigit(char):
			hasNum.valid = true
		case unicode.IsLower(char):
			hasLowerChar.valid = true
		case unicode.IsUpper(char):
			hasUpperChar.valid = true
		case !unicode.IsDigit(char) && !unicode.IsLetter(char):
			hasSymbol.valid = true
		}
	}

	switch {
	case !hasNum.valid:
		return false, hasNum.message
	case !hasUpperChar.valid:
		return false, hasUpperChar.message
	case !hasLowerChar.valid:
		return false, hasLowerChar.message
	case !hasSymbol.valid:
		return false, hasSymbol.message
	}
	return true, ""
}

type registerBody struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (app App) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte("Error when reading the body of the request"))
		return
	}

	data := registerBody{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.Write([]byte("Invalid json"))
		return
	}

	valid, message := IsValidPassword(&data.Password)
	if !valid {
		w.Write([]byte(message))
		return
	}
	if data.Nickname == "" || data.Email == "" {
		w.Write([]byte("Nickname and email are required"))
		return
	}

	hashedPassword, err := hashPassword(&data.Password)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = app.DB.CreateUser(&data.Nickname, &data.Email, &hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("Created user"))
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
		res, error := jsonError(err.Error(), 500)
		if error != nil {
			http.Error(w, error.Error(), 500)
		}
		w.Write(res)
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

	isSamePass := verifyPassword(&data.Password, &user.Password)
	if !isSamePass {
		http.Error(w, "Incorrect email or password", http.StatusNotFound)
		return
	}

	strUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(strUser)
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
