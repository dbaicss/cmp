package model

type ListPod struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	PodIP     string `json:"pod_ip"`
	HostIp    string `json:"host_ip"`
	StartTime string `json:"start_time"`
}
