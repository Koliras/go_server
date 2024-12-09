package api

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 14)
	return string(bytes), err
}

func verifyPassword(password, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*password))
	return err == nil
}

type JsonError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func jsonError(message string, code int) ([]byte, error) {
	jsonErr := JsonError{message, code}
	strError, err := json.Marshal(jsonErr)
	if err != nil {
		return nil, err
	}
	return strError, nil
}
