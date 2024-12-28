package middleware

import (
	"net/http"

	database "github.com/Koliras/go_server/db"
	"github.com/Koliras/go_server/templ"
	"github.com/golang-jwt/jwt/v5"
)

const JwtTokenCookieName = "Go-Server"

// WARN: change later to real secret
var JwtKey = []byte("some_secret_key")

type JwtClaims struct {
	Email    string
	Nickname string
	jwt.RegisteredClaims
}

func JwtAuth(next func(http.ResponseWriter, *http.Request), db database.DbInstance) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(JwtTokenCookieName)
		if err != nil {
			templ.HtmlAuthError(w, "Could not find authentication cookie")
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			templ.HtmlAuthError(w, "Could not parse authentication token")
			return
		}

		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			templ.HtmlAuthError(w, "Incorrect authentication token")
			return
		}
		exists, err := db.UserExists(&claims.Email)
		if err != nil {
			templ.HtmlInternalError(w, "Internal server error while authenticating")
			return
		} else if !exists {
			templ.HtmlAuthError(w, "User with such email does not exist")
			return
		}

		next(w, r)
	}
}
