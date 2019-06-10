package logic

import (
	"cmp-server/dal/db"
	"cmp-server/model"
)

//查询apply_cvm记录
func GetResourceList() (count []*model.ResourceData, err error) {
	count, err = db.GetResourceList()
	if err != nil {
		return nil, err
	}
	return
}
