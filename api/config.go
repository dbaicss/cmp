package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//单例对象
var (
	G_config *Config
)

//程序配置
type Config struct {
	MysqlAddr       string `json:"mysqlAddr"`
	MysqlPort       int    `json:"mysqlPort"`
	MysqlUser       string `json:"mysqlUser"`
	MysqlPassword   string `json:"mysqlPassword"`
	MysqlDatabase   string `json:"mysqlDatabase"`
	MysqlProtocol   string `json:"mysqlProtocol"`
	WebRoot         string `json:"webroot"`
	ApiAddr         string `json:"apiAddr"`
	ApiPort         int    `json:"apiPort"`
	ApiReadTimeout  int    `json:"apiReadTimeout"`
	ApiWriteTimeout int    `json:"apiWriteTimeout"`
	KubeConfig      string `json:"kubeconfig"`
	K8sConfig       string `json:"k8s_config"`
	SecretId        string `json:"secretId"`
	SecretKey       string `json:"secretKey"`
	Region          string `json:"region"`
}

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	//加载配置
	if content, err = ioutil.ReadFile(filename); err != nil {
		fmt.Println("读取配置失败", err)
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		fmt.Println("配置反序列化失败:", err)
		return
	}

	//单例初始化
	G_config = &conf
	return
}
