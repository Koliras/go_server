package api

import (
	"database/sql"
	"net/http"
)

func Routes(con *sql.DB) http.Handler {
	app := App{
		DbInstance{con},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.Healthcheck)
	mux.HandleFunc("POST /auth/register", app.Register)
	mux.HandleFunc("GET /users", app.GetAllUsers)

	return mux
}
