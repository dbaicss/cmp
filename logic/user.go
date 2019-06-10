package logic

import (
	"net/http"

	"cmp-server/auth"
	c "cmp-server/common"
	"cmp-server/dal/db"
	"cmp-server/model"

	"fmt"
	"encoding/json"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		un,ps string
		err error
		)
	//1.解析post表单
	if err = r.ParseForm(); err != nil {
		fmt.Println("parse form err:",err)
		return
	}
	un = r.PostForm.Get("username")
	ps = r.PostForm.Get("password")
	user.UserName = un
	user.Password = ps
	if err != nil || un == "" || ps == "" {
		c.ResponseWithJson(w, http.StatusBadRequest,
			c.Response{ErrNo: http.StatusBadRequest, Msg: "bad params"})
		return
	}
	err = db.Insert(user)
	if err != nil {
		c.ResponseWithJson(w, http.StatusInternalServerError,
			c.Response{ErrNo: http.StatusInternalServerError, Msg: "internal error"})
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		err error
	)

	err = json.NewDecoder(r.Body).Decode(&user)
	if  err != nil {
		c.ResponseWithJson(w, http.StatusBadRequest,
			c.Response{ErrNo: http.StatusBadRequest, Msg: "bad params"})
		return
	}
	exist := db.IsExist(user.UserName)
	if exist {
		token, _ := auth.GenerateToken(&user)
		c.ResponseWithJson(w, http.StatusOK,
			c.Response{ErrNo: http.StatusOK, Data: model.JwtToken{Token: token}})
	} else {
		c.ResponseWithJson(w, http.StatusNotFound,
			c.Response{ErrNo: http.StatusNotFound, Msg: "the user not exist"})
	}
}
