package api

type User struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password Pass   `json:"password"`
}
