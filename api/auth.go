package api

import (
	"encoding/json"
	"io"
	"net/http"
	"unicode"
)

type Pass string

func (p Pass) IsValid() (bool, string) {
	if len(p) < 8 {
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
	for _, char := range p {
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
	Password Pass   `json:"password"`
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

	valid, message := data.Password.IsValid()
	if !valid {
		w.Write([]byte(message))
		return
	}
	if data.Nickname == "" || data.Email == "" {
		w.Write([]byte("Nickname and email are required"))
	}

	err = app.DB.CreateUser(data.Nickname, data.Email, string(data.Password))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte("Created user"))
}
