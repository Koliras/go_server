package api

import (
	"database/sql"
	"net/http"

	"github.com/Koliras/go_server/db"
	"github.com/Koliras/go_server/middleware"
)

type App struct {
	DB db.DbInstance
}

func Routes(con *sql.DB) http.Handler {
	app := App{
		db.DbInstance{con},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", app.Healthcheck)
	mux.HandleFunc("POST /auth/login", middleware.JwtAuth(app.Login, app.DB))
	mux.HandleFunc("POST /auth/register", app.Register)
	mux.HandleFunc("GET /users", app.GetAllUsers)

	return mux
}
