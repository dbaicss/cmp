package common

import (
	"encoding/json"
	"net/http"
)

//http接口应答
type Response struct {
	ErrNo int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response *Response
	)
	response = &Response{
		ErrNo: errno,
		Msg:   msg,
		Data:  data,
	}
	if resp, err = json.Marshal(response); err != nil {
		return
	}
	return
}

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}