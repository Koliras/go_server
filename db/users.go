package db

type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db DbInstance) CreateUser(nickname, email, password string) error {
	stmnt := "INSERT INTO users (nickname, email, password) VALUES(?, ?, ?)"
	_, err := db.Exec(stmnt, nickname, email, password)
	return err
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
