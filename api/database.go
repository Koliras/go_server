package api

import "database/sql"

type DbInstance struct {
	*sql.DB
}
type App struct {
	DB DbInstance
}

func (db DbInstance) CreateUser(nickname, email, password string) error {
	stmnt := "INSERT INTO users (nickname, email, password) VALUES(?, ?, ?)"
	_, err := db.Exec(stmnt, nickname, email, password)
	return err
}

// func (DbModel) AllUsers() ([]User, error) {
// 	stmt := "SELECT * FROM users"
// }
