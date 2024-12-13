package api

import (
	"encoding/json"

	"golang.org/x/crypto/bcrypt"
)

func hashString(str *string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(*str), 14)
	return string(bytes), err
}

func compareStrWithHash(str, hash *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*hash), []byte(*str))
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
