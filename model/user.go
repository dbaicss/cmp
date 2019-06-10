package model

type User struct {
	UserName string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type JwtToken struct {
	Token string `json:"token"`
}
