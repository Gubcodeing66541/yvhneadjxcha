package Response

import (
	"server/App/Constant"
)

type RespDomainList struct {
	Id                   int                 `json:"id"`
	Domain               string              `json:"domain"`
	Type                 string              `json:"type"`
	BindServiceId        int                 `json:"bind_service_id"`
	WeChatBanStatus      string              `json:"we_chat_ban_status"`
	ServiceId            int                 `json:"service_id"`
	Name                 string              `json:"name"`
	Status               string              `json:"status"`
	Username             string              `json:"username"`
	CreateTime           Constant.SystemTime `json:"create_time"`
	ServiceManagerName   string              `json:"service_manager_name"`
	ServiceManagerId     int                 `json:"service_manager_id"`
	ServiceManagerMember string              `json:"service_manager_member"`
}
