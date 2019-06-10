package db

import (
	"cmp-server/model"
)
//查询asset
func GetAssetList() (crontabList []*model.AssetInfo, err error) {

	sqlstr := `select id, model, sn, location, ip, port, kind, status  from asset`

	err = DB.Select(&crontabList, sqlstr)
	return
}
