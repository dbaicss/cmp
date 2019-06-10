package api

import (
	"net/http"
	"strings"
	"time"

	c "cmp-server/common"
	"cmp-server/logic"
	"cmp-server/model"

	"github.com/gorilla/websocket"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	yj "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/yunjing/v20180228"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	v1 "k8s.io/api/core/v1"

	"cmp-server/api/ws"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

)

var (
	clientset *kubernetes.Clientset
)

// ssh流式处理器
type streamHandler struct {
	wsConn      *ws.WsConnection
	resizeEvent chan remotecommand.TerminalSize
}

// web终端发来的包
type xtermMessage struct {
	MsgType string `json:"type"`  // 类型:resize客户端调整终端, input客户端输入
	Input   string `json:"input"` // msgtype=input情况下使用
	Rows    uint16 `json:"rows"`  // msgtype=resize情况下使用
	Cols    uint16 `json:"cols"`  // msgtype=resize情况下使用
}

// executor回调获取web是否resize
func (handler *streamHandler) Next() (size *remotecommand.TerminalSize) {
	ret := <-handler.resizeEvent
	size = &ret
	return
}

// executor回调读取web端的输入
func (handler *streamHandler) Read(p []byte) (size int, err error) {
	var (
		msg      *ws.WsMessage
		xtermMsg xtermMessage
	)

	// 读web发来的输入
	if msg, err = handler.wsConn.WsRead(); err != nil {
		fmt.Println("web发来的输入错误:", err)
		return
	}

	// 解析客户端请求
	if err = json.Unmarshal(msg.Data, &xtermMsg); err != nil {
		fmt.Println("解析客户端请求错误:", err)
		return
	}

	//web ssh调整了终端大小
	if xtermMsg.MsgType == "resize" {
		// 放到channel里，等remotecommand executor调用我们的Next取走
		handler.resizeEvent <- remotecommand.TerminalSize{Width: xtermMsg.Cols, Height: xtermMsg.Rows}
	} else if xtermMsg.MsgType == "input" { // web ssh终端输入了字符
		// copy到p数组中
		size = len(xtermMsg.Input)
		//copy(p, xtermMsg.Input)
		copy(p, xtermMsg.Input)
	}
	return
}

// executor回调向web端输出
func (handler *streamHandler) Write(p []byte) (size int, err error) {
	size = len(p)
	err = handler.wsConn.WsWrite(websocket.TextMessage, p)
	return
}

func WsHandler(resp http.ResponseWriter, req *http.Request) {
	var (
		wsConn        *ws.WsConnection
		restConf      *rest.Config
		sshReq        *rest.Request
		podName       string
		podNs         string
		containerName string
		executor      remotecommand.Executor
		handler       *streamHandler
		err           error
	)

	// 解析GET参数
	if err = req.ParseForm(); err != nil {
		return
	}
	podNs = req.Form.Get("podNs")
	podName = req.Form.Get("podName")
	containerName = req.Form.Get("containerName")

	// 得到websocket长连接
	if wsConn, err = ws.InitWebsocket(resp, req); err != nil {
		fmt.Println("解析kube config失败:", err)
		return
	}

	// 获取pods

	// 获取k8s rest client配置
	if restConf, err = GetRestConf(); err != nil {
		fmt.Println("生成restConf失败,err:", err)
		goto END
	}

	// 生成clientset配置
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		fmt.Println("生成clientset失败,err:", err)
		goto END
	}
	// URL长相:
	// https://172.18.11.25:6443/api/v1/namespaces/default/pods/nginx-deployment-5cbd8757f-d5qvx/exec?command=sh&container=nginx&stderr=true&stdin=true&stdout=true&tty=true
	sshReq = clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(podNs).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// 创建到容器的连接
	if executor, err = remotecommand.NewSPDYExecutor(restConf, "POST", sshReq.URL()); err != nil {
		fmt.Println("执行ssh new spdy err:", err)
		goto END
	}

	// 配置与容器之间的数据流处理回调
	handler = &streamHandler{wsConn: wsConn, resizeEvent: make(chan remotecommand.TerminalSize)}
	if err = executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		fmt.Println("处理容器数据流回调失败:", err)
		goto END
	}
	return

END:
	fmt.Println(err)
	wsConn.WsClose()
}

// 获取k8s restful client配置
func GetRestConf() (restConf *rest.Config, err error) {
	var (
		kubeconfig []byte
	)

	// 读kubeconfig文件
	if kubeconfig, err = ioutil.ReadFile(G_config.KubeConfig); err != nil {
		goto END
	}
	// 生成rest client配置
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		goto END
	}
END:
	return
}

func HandleLogList(resp http.ResponseWriter, req *http.Request) {
	var (
		tailLines     int64 = 1000
		cont          *v1.Pod
		contain       []v1.Container
		restConf      *rest.Config
		res           rest.Result
		podName       string
		podNs         string
		r             *rest.Request
		containerName string
		respon        []byte
		logs          []byte
		err           error
	)
	// 获取k8s rest client配置
	if restConf, err = GetRestConf(); err != nil {
		fmt.Println("生成restConf失败,err:", err)
		goto END
	}

	// 生成clientset配置
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		fmt.Println("生成clientset失败,err:", err)
		goto END
	}
	// 解析GET参数
	if err = req.ParseForm(); err != nil {
		return
	}
	podNs = req.URL.Query().Get("ns")
	podName = req.URL.Query().Get("name")
	cont, err = clientset.CoreV1().Pods(podNs).Get(podName, metav1.GetOptions{})
	if err != nil {
		fmt.Println("get container failed,err:", err)
		goto END
	}
	contain = cont.Spec.Containers
	containerName = contain[0].Name
	fmt.Println("podNs:", podNs, "podName:", podName, "containerName:", containerName)
	r = clientset.CoreV1().Pods(podNs).GetLogs(podName, &v1.PodLogOptions{Container: containerName, TailLines: &tailLines})

	// 发送请求
	if res = r.Do(); res.Error() != nil {
		err = res.Error()
		fmt.Println("do:", err)
		goto END
	}

	// 获取结果
	if logs, err = res.Raw(); err != nil {
		fmt.Println("raw:", err)
		goto END
	}

	if respon, err = c.BuildResponse(0, "succ", string(logs)); err == nil {
		resp.Write(respon)
	}
	return
END:
	fmt.Println(err)
	//异常应答
	if respon, err = c.BuildResponse(-1, err.Error(), []string{}); err == nil {
		resp.Write(respon)
	}
}

//查询crontab
func HandleCrontabList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*model.CrontabInfo
		respon  []byte
		err     error
	)
	//获取任务列表数据
	if jobList, err = logic.GetCrontabList(); err != nil {
		goto Err
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", jobList); err == nil {
		resp.Write(respon)
	}
	return
Err:
	//异常应答
	if respon, err = c.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(respon)
	}
}

//查询asset
func HandleAssetList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*model.AssetInfo
		respon  []byte
		err     error
	)
	//获取任务列表数据
	if jobList, err = logic.GetAssetList(); err != nil {
		goto Err
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", jobList); err == nil {
		resp.Write(respon)
	}
	return
Err:
	//异常应答
	if respon, err = c.BuildResponse(-1, err.Error(), []string{}); err == nil {
		resp.Write(respon)
	}
}

//查询asset server
func HandleServerList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*model.ServerInfo
		respon  []byte
		err     error
	)
	//获取任务列表数据
	if jobList, err = logic.GetServerList(); err != nil {
		goto Err
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", jobList); err == nil {
		resp.Write(respon)
	}
	return
Err:
	//异常应答
	if respon, err = c.BuildResponse(-1, err.Error(), []string{}); err == nil {
		resp.Write(respon)
	}
}

//查询apply_cvm表,申请资源记录
func HandleAuditList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*model.ResourceData
		respon  []byte
		err     error
	)
	//获取任务列表数据
	if jobList, err = logic.GetResourceList(); err != nil {
		goto Err
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", jobList); err == nil {
		resp.Write(respon)
	}
	return
Err:
	//异常应答
	if respon, err = c.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(respon)
	}
}

//查询k8s所有的namespace
func HandleNameSpaceList(resp http.ResponseWriter, req *http.Request) {
	var (
		ns_names []string
		nss      *v1.NamespaceList
		respon   []byte
		err      error
	)
	config, err := clientcmd.BuildConfigFromFlags("", G_config.KubeConfig)
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)

	nss, err = clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	for _, i := range nss.Items {
		ns_names = append(ns_names, i.Name)
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", ns_names); err == nil {
		resp.Write(respon)
	}
	return
}

//查询k8s指定namespace下面pod列表
func HandlePodList(resp http.ResponseWriter, req *http.Request) {
	var (
		respon []byte
		name   string
		pods   *v1.PodList
		res    []model.ListPod
		err    error
	)
	//config, err := rest.InClusterConfig()
	config, err := clientcmd.BuildConfigFromFlags("", G_config.KubeConfig)
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	//1.解析get请求参数
	name = req.URL.Query().Get("ns")
	pods, err = clientSet.CoreV1().Pods(name).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, i := range pods.Items {
		var pod = model.ListPod{}
		pod.Name = i.Name
		pod.Namespace = i.Namespace
		pod.Status = string(i.Status.Phase)
		pod.PodIP = i.Status.PodIP
		now := time.Unix(i.Status.StartTime.Time.Unix(), 0).Format("2006-01-02 03:04:05 PM")
		pod.StartTime = now
		pod.HostIp = i.Status.HostIP
		res = append(res, pod)
	}

	if len(res) == 0 {
		if respon, err = c.BuildResponse(0, "succ", []string{}); err == nil {
			resp.Write(respon)
		}
	} else {
		//正常应答
		if respon, err = c.BuildResponse(0, "succ", res); err == nil {
			resp.Write(respon)
		}
	}
	return

}

type Pro struct {
	CreateTime  string `json:"createtime"`
	FullPath    string `json:"fullpath"`
	//Id          int    `json:"id"`
	//MachineIp   string `json:"machineip"`
	//MachineName string `json:"machinename"`
	Pid         int    `json:"pid"`
	//Platform    string `json:"platform"`
	Ppid        int    `json:"ppid"`
	ProcessName string `json:"processname"`
	Username    string `json:"username"`
	//Uuid        string `json:"uuid"`
}

type Ret struct {
	Processes  []Pro
	RequestId  string
	TotalCount int
}

type Res struct {
	Response Ret
}

//查询cvm上服务列表信息
func HandleAssetServiceInfoList(resp http.ResponseWriter, req *http.Request) {
	credential := common.NewCredential(G_config.SecretId, G_config.SecretKey)

	// 实例化一个客户端配置对象
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "yunjing.tencentcloudapi.com"
	client, _ := yj.NewClient(credential, G_config.Region, cpf)
	request := yj.NewDescribeProcessesRequest()
	//返回数据的条数
	request.Limit = common.Uint64Ptr(100)
	request.Uuid = common.StringPtr("ff442974-d0fa-11e7-8f8b-98be94219792")
	response, err := client.DescribeProcesses(request)
	// 处理异常
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	var res Res
	err = json.Unmarshal([]byte(response.ToJsonString()), &res)
	if err != nil {
		log.Println("response unmarshal failed,err:", err)
	}
	var retu []Pro
	var respon []byte
	for _, i := range res.Response.Processes {
		if strings.HasPrefix(i.ProcessName, "ssh") || strings.HasPrefix(i.ProcessName, "sudo") || strings.HasPrefix(i.ProcessName, "su") || strings.HasPrefix(i.ProcessName, "kworker") || strings.HasPrefix(i.ProcessName, "bash") || strings.HasPrefix(i.ProcessName, "md") || strings.HasPrefix(i.ProcessName, "kworker") || strings.HasPrefix(i.ProcessName, "bash") || strings.HasPrefix(i.ProcessName, "khub") || strings.HasPrefix(i.ProcessName, "kworker") || strings.HasPrefix(i.ProcessName, "bash") || strings.HasPrefix(i.ProcessName, "dev") || strings.HasPrefix(i.ProcessName, "khungtaskd") {
			continue
		}
		var pro = Pro{}
		pro.CreateTime = i.CreateTime
		pro.FullPath = i.FullPath
		pro.Pid = i.Pid
		pro.Ppid = i.Ppid
		pro.ProcessName = i.ProcessName
		pro.Username = i.Username
		retu = append(retu, pro)
	}
	//正常应答
	if respon, err = c.BuildResponse(0, "succ", retu); err == nil {
		resp.Write(respon)
	}else {
		if respon, err = c.BuildResponse(-1, err.Error(), nil); err == nil {
			resp.Write(respon)
		}
	}
	return
}
