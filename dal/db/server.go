package db

import (
	"cmp-server/model"
)
//查询asset
func GetServerList() (crontabList []*model.ServerInfo, err error) {

	sqlstr := `select id, name, location, model, sn,  ip, purch_time, used, department, audit, isvirtual, status, description  from server`

	err = DB.Select(&crontabList, sqlstr)
	return
}
