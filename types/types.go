package types

type MyTask struct {
	TaskCpu     float64           `json:"task_cpu"`
	TaskMem     float64           `json:"task_mem"`
	ID          string            `json:"id"`
	Cmd         string            `json:"cmd"`
	Env         map[string]string `json:"env"`
	Image       string            `json:"image"`
	SlaveId     string            `json:"slave_id"`
	Hostname    string            `json:"hostname"`
	Name        string            `json:"name"`
	FrameworkId string            `json:"framework_id"`
	Status      int               `json:"status"`
	Count       int               `json:"count"`
}

type TaskRequest struct {
	Cpu      float64           `json:"cpu"`
	Mem      float64           `json:"mem"`
	Cmd      string            `json:"cmd"`
	Env      map[string]string `json:"env"`
	Name     string            `json:"name"`
	Image    string            `json:"image"`
	Count    int               `json:"count"`
	Hostname string            `json:"hostname"`
}

type SlaveNode struct {
	Hostname   string `json:"hostname"`
	Attachment string `json:"attachment"`
}

type Config struct {
	Master  string
	Address string
	Port    uint16
	Webui   string
	Name    string
	User    string
}
