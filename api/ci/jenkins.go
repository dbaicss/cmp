package ci
//package main


import (
	"github.com/bndr/gojenkins"
	"fmt"
	"crypto/tls"
	"net/http"
)

func BuildExistJob(jobname string)  {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport:tr}

	//初始化jenkins

	jenkins,err := gojenkins.CreateJenkins(client, "https://xxxx.cn", "admin", "xxxxxx").Init()

	if err != nil {
		panic("Init Jenkins wrong...")
	}
	//job, err := jenkins.GetJob("account-auth")
	job, err := jenkins.GetJob(jobname)
	if err != nil {
		panic("Job does not exist")
	}

	param,err := job.GetParameters()
	//获取parameters,一般都是tag发布的版本号
	res := make(map[string]string)
	for _,p := range param {
		fmt.Printf("params:%#v\n",p.DefaultParameterValue)
		res[p.DefaultParameterValue.Name] = p.DefaultParameterValue.Value.(string)
	}
	//buildId,err := jenkins.BuildJob("account-auth",res)
	//data, err := jenkins.GetBuild("account-auth", buildId)

	buildId,err := jenkins.BuildJob(jobname,res)
	data, err := jenkins.GetBuild(jobname, buildId)
	if err != nil {
		panic(err)
	}

	if "SUCCESS" == data.GetResult() {
		fmt.Println("This build succeeded")
	}
}

func BuildNewJob()  {
	fmt.Println("ok")
}

//func main()  {
//	BuildExistJob("icx-blog")
//}