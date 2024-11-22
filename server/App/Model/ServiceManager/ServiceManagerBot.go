package ServiceManager

type ServiceManagerBot struct {
	Id               int
	Status           string `json:"status"` // stop停止 run启动
	ServiceManagerId int    `json:"service_manager_id"`
	Name             string `json:"name"` // 默认就叫智能机器人
	Head             string `json:"head"`
	Hello            string `json:"hello"`
}
