package model

type ResourceData struct {
	Id            int    `json:"id" db:"id"`
	Kind          string `json:"kind" db:"kind"`
	Disk          string `json:"disk" db:"disk"`
	DiskType      string `json:"disk_type" db:"disk_type"`
	OsVersion     string `json:"os_version" db:"os_version"`
	Network       int    `json:"network" db:"network"`
	InstanceCount int    `json:"instance_count" db:"instance_count"`
	LiveTime      int    `json:"live_time" db:"live_time"`
	DeployTime    string `json:"deploy_time" db:"deploy_time"`
	Description   string `json:"description" db:"description"`
	ApplyName     string `json:"apply_name" db:"apply_name"`
	SpStatus      int    `json:"sp_status" db:"sp_status"`
	ApplyTime     int64  `json:"apply_time" db:"apply_time"`
	SpNum         int64  `json:"sp_num" db:"sp_num"`
	Spname        string `json:"spname" db:"sp_name"`
	ApplyOrg      string `json:"apply_org" db:"apply_org"`
	ApprovalName  string `json:"approval_name" db:"approval_name"`
	NotifyName    string `json:"notify_name" db:"notify_name"`
	ApplyUserId   string `json:"apply_user_id" db:"apply_user_id"`
}
