package db

import (
	"cmp-server/model"
	"fmt"
	"database/sql"
)

func IsExist(username string) bool {
	var (
		query string
		exists bool
	)
	sqlstr := `select id from user where username=?`
	query = fmt.Sprintf("SELECT exists (%s)", sqlstr)
	err :=DB.QueryRow(query,username).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false
	}
	return true
}

func Insert(user  model.User) (err error)  {
	sqlstr := `insert into user(username,password) values(?,?)`
	_, err = DB.Exec(sqlstr,user.UserName,user.Password)
	return
}