package logic

import (
	"cmp-server/dal/db"
	"cmp-server/model"
)
//查询crontab记录
func GetAssetList() (crontabList []*model.AssetInfo, err error) {
	crontabList, err = db.GetAssetList()
	if err != nil {
		return nil, err
	}
	return
}
