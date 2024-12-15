package middleware

import (
	"net/http"

	database "github.com/Koliras/go_server/db"
	"github.com/golang-jwt/jwt/v5"
)

const JwtTokenCookieName = "Go-Server"

// WARN: change later to real secret
var JwtKey = []byte("some_secret_key")

type JwtClaims struct {
	Email    string `json:"email"`
	Nickname string `json:"nickcname"`
	jwt.RegisteredClaims
}

func JwtAuth(next func(http.ResponseWriter, *http.Request), db database.DbInstance) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(JwtTokenCookieName)
		if err != nil {
			http.Error(w, "No auth cookie found", http.StatusUnauthorized)
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			http.Error(w, "Error when parsing authentication token", 500)
			return
		}

		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			http.Error(w, "Incorrect authentication token after parsing", 500)
			return
		}
		exists, err := db.UserExists(&claims.Email)
		if err != nil {
			http.Error(w, "Database error when checking if user exists", 500)
			return
		} else if !exists {
			http.Error(w, "User does not exist", 500)
			return
		}

		next(w, r)
	}
}
