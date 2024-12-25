package db

import "database/sql"

type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db DbInstance) CreateUser(nickname, email, password *string) error {
	stmnt := "INSERT INTO users (nickname, email, password) VALUES(?, ?, ?)"
	_, err := db.Exec(stmnt, nickname, email, password)
	return err
}

func (db DbInstance) GetUserByEmail(email *string) (User, error) {
	stmt := "SELECT id, nickname, email, password FROM users WHERE email=? LIMIT 1"
	u := User{}
	err := db.QueryRow(stmt, email).Scan(&u.Id, &u.Nickname, &u.Email, &u.Password)
	return u, err
}

func (db DbInstance) UserExists(email *string) (bool, error) {
	if len(*email) == 0 {
		return false, nil
	}
	stmt := "SELECT 1 FROM users WHERE email=? LIMIT 1"
	u := 0
	err := db.QueryRow(stmt, email).Scan(&u)

	if err == sql.ErrNoRows {
		return false, nil
	}
	return u == 1, err
}

func (db DbInstance) NicknameUsed(nickname *string) (bool, error) {
	if len(*nickname) == 0 {
		return false, nil
	}
	stmt := "SELECT 1 FROM users WHERE nickname=? LIMIT 1"
	u := 0
	err := db.QueryRow(stmt, nickname).Scan(&u)

	if err == sql.ErrNoRows {
		return false, nil
	}
	return u == 1, err
}

func (db DbInstance) AllUsers() ([]User, error) {
	stmt := "SELECT id, nickname, email, password FROM users ORDER BY id ASC"
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}

	users := []User{}
	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Id, &u.Nickname, &u.Email, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
