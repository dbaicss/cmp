package main



import (
	"net/http"
	"fmt"
	"flag"
	"runtime"
	"cmp-server/dal/db"
	"strconv"
	"cmp-server/routes"
	"cmp-server/api"
)

var (
	confFile  string //配置文件路径

)

//初始化线程
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	flag.StringVar(&confFile, "config", "./api.json", "config path")
	flag.Parse()
}

func main() {

	var (
		err error
		dns string
	)
	r := routes.NewRouter()
	//初始化命令行参数
	initArgs()

	//初始化线程
	initEnv()

	//初始化日志
	if err = api.InitLog(); err != nil {
		fmt.Println("日誌初始化異常:", err)
		goto Err
	}

	//加载配置
	if err = api.InitConfig(confFile); err != nil {
		fmt.Println("配置初始化異常:", err)
		goto Err
	}

	//root:@tcp(127.0.0.1:3306)/bloger?parseTime=true
	//初始化mysql配置
	dns = api.G_config.MysqlUser + ":" + api.G_config.MysqlPassword + "@" + api.G_config.MysqlProtocol + "(" + api.G_config.MysqlAddr + ":" + strconv.Itoa(api.G_config.MysqlPort) + ")" + "/" + api.G_config.MysqlDatabase + "?parseTime=true"
	err = db.Init(dns)
	if err != nil {
		fmt.Println("exec mysql error:",err)
		panic(err)
	}

	http.ListenAndServe(fmt.Sprintf("%s:%d", api.G_config.ApiAddr, api.G_config.ApiPort), r)

Err:
	fmt.Println("加载配置异常:", err)
}