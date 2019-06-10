package model

type AssetInfo struct {
	Id                 int `json:"instanceId" db:"id"`
	Model      string `json:"model" db:"model"`
	Sn        string `json:"sn" db:"sn"`
	Location        string `json:"location" db:"location"`
	Ip  string `json:"ip" db:"ip"`
	Port int `json:"port" db:"port"`
	Kind  string `json:"kind" db:"kind"`
	Status  string     `json:"status" db:"status"`
}
