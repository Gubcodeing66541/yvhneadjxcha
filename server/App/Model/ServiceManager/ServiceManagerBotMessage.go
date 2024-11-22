package ServiceManager

import "time"

type ServiceManagerBotMessage struct {
	Id               int       `json:"id" uri:"id"  form:"id" `
	ServiceManagerId int       `json:"service_manager_id"`
	Problem          string    `json:"problem" uri:"problem"  form:"problem"`
	Answer           string    `json:"answer" uri:"answer"  form:"answer"`
	CreateTime       time.Time `json:"create_time"`
}
