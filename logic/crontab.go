package logic

import (
	"cmp-server/dal/db"
	"cmp-server/model"
)

//查询crontab记录
func GetCrontabList() (crontabList []*model.CrontabInfo, err error) {
	crontabList, err = db.GetCrontabList()
	if err != nil {
		return nil, err
	}
	return
}
