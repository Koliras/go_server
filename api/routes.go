package api

import (
	"database/sql"
	"net/http"

	"github.com/Koliras/go_server/db"
)

type App struct {
	DB db.DbInstance
}

func Routes(con *sql.DB) http.Handler {
	app := App{
		db.DbInstance{con},
	}
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /register", app.StaticRegisterForm)
	mux.HandleFunc("GET /login", app.StaticLoginForm)

	mux.HandleFunc("GET /api/healthcheck", app.Healthcheck)
	mux.HandleFunc("POST /api/auth/login", app.Login)
	mux.HandleFunc("POST /api/auth/register", app.Register)
	mux.HandleFunc("GET /api/users", app.GetAllUsers)

	return mux
}
