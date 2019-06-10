package db

import (
	"cmp-server/model"
)

//查询crontab
func GetCrontabList() (crontabList []*model.CrontabInfo, err error) {

	sqlstr := `select 
						id, instanceState, createdTime, expiredTime, publicIpAddresses, privateIpAddresses, imageId, instanceType, instanceName, osName,
					zone, resource, internetAccessible
					from 
						crontab`

	err = DB.Select(&crontabList, sqlstr)
	return
}
