package model

type ServerInfo struct {
	Id                 int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Location        string `json:"location" db:"location"`
	Model      string `json:"model" db:"model"`
	Sn        string `json:"sn" db:"sn"`
	Ip  string `json:"ip" db:"ip"`
	PurchTime  string  `json:"purch_time" db:"purch_time"`
	Used string `json:"used" db:"used"`
	Department string `json:"department" db:"department"`
	Audit string `json:"audit" db:"audit"`
	Isvirtual int  `json:"isvirtual" db:"isvirtual"`
	Status  string     `json:"status" db:"status"`
	Description  string `json:"description" db:"description"`
}
