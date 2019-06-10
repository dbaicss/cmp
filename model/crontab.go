package model

type CrontabInfo struct {
	Id                 string `json:"instanceId" db:"id"`
	InstanceState      string `json:"instanceState" db:"instanceState"`
	CreatedTime        string `json:"createdTime" db:"createdTime"`
	ExpiredTime        string `json:"expiredTime" db:"expiredTime"`
	PublicIpAddresses  string `json:"publicIpAddresses" db:"publicIpAddresses"`
	PrivateIpAddresses string `json:"privateIpAddresses" db:"privateIpAddresses"`
	ImageId            string `json:"imageId" db:"imageId"`
	InstanceType       string `json:"instanceType" db:"instanceType"`
	InstanceName       string `json:"instanceName" db:"instanceName"`
	OsName             string `json:"osName" db:"osName"`
	Zone               string `json:"zone" db:"zone"`
	Resource           string `json:"resource" db:"resource"`
	InternetAccessible int    `json:"internetAccessible" db:"internetAccessible"`
}
