package utils

import (
	"net/mail"
	"unicode"

	"github.com/Koliras/go_server/templ"
)

type RegisterInput struct {
	Nickname       string `json:"nickname"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`
	Email          string `json:"email"`
}

func ValidateRegisterData(d RegisterInput) templ.RegisterFormData {
	r := templ.RegisterFormData{}

	r.PasswordErrors = passErrors(&d.Password)
	if d.Password != d.RepeatPassword {
		r.RepeatPasswordErrors = append(r.RepeatPasswordErrors, "Passwords are not the same")
	}
	if _, err := mail.ParseAddress(d.Email); err != nil {
		if len(d.Email) == 0 {
			r.EmailErrors = append(r.EmailErrors, "Email is required")
		} else {
			r.EmailErrors = append(r.EmailErrors, "Invalid email")
		}
	}
	if len(d.Nickname) > 20 {
		r.NicknameErrors = append(r.NicknameErrors, "Nickname is too long. Max length is 20 characters")
	} else if len(d.Nickname) == 0 {
		r.NicknameErrors = append(r.NicknameErrors, "Nickname is required")
	}

	return r
}

func passErrors(p *string) []string {
	var hasNum, hasUpperChar, hasLowerChar, hasSymbol bool
	r := make([]string, 0, 2)

	if len(*p) < 8 {
		r = append(r, "Password is too short")
	}

	for _, char := range *p {
		switch {
		case unicode.IsDigit(char):
			hasNum = true
		case unicode.IsLower(char):
			hasLowerChar = true
		case unicode.IsUpper(char):
			hasUpperChar = true
		case !unicode.IsDigit(char) && !unicode.IsLetter(char):
			hasSymbol = true
		}
	}

	if !hasNum || !hasUpperChar || !hasLowerChar || !hasSymbol {
		r = append(r, "Password has to contain at least 1 number, uppercase letter, lowercase letter and symbol")
	}
	return r
}
