package db

import "cmp-server/model"

//查询apply_cvm表
func GetResourceList() (rdataList []*model.ResourceData, err error) {

	sqlstr := `select kind,disk,disk_type,os_version,network,instance_count,live_time,
				deploy_time,description,apply_name,sp_status,apply_time,sp_num,sp_name,apply_org,approval_name,notify_name,apply_user_id from  apply_cvm order by apply_time desc`
	err = DB.Select(&rdataList, sqlstr)
	return
}
