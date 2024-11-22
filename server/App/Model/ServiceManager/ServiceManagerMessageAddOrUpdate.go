package ServiceManager

import (
	"server/App/Constant"
)

type ServiceManagerMessage struct {
	Id               int                 `json:"id"`
	Title            string              `json:"title"`
	Content          string              `json:"content"`
	Type             string              `json:"type"`
	ServiceManagerId int                 `json:"service_manager_id"`
	AddServiceId     int                 `json:"add_service_id"`
	Status           string              `json:"status"` // 表示  apply 审批中 agree 同意 refuse 拒绝
	CreateTime       Constant.SystemTime `json:"create_time" `
	UpdateTime       Constant.SystemTime `json:"update_time"`
}
