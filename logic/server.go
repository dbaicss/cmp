package logic

import (
	"cmp-server/dal/db"
	"cmp-server/model"
)
//查询crontab记录
func GetServerList() (crontabList []*model.ServerInfo, err error) {
	crontabList, err = db.GetServerList()
	if err != nil {
		return nil, err
	}
	return
}
