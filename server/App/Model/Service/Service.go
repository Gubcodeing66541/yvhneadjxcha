package Service

import "time"

// 客服表
type Service struct {
	Id               int       `json:"id" gorm:"primary_key;AUTO_INCREMENT" `
	ServiceId        int       `json:"service_id"`
	ServiceManagerId int       `json:"service_manager_id"`
	Name             string    `json:"name"`
	Head             string    `json:"head"`
	Mobile           string    `json:"mobile"`
	Username         string    `json:"username"`
	Type             string    `json:"type"` //auth push
	Role             string    `json:"role"`
	Day              int       `json:"day"`
	Code             string    `json:"code"`
	IsActivate       int       `json:"is_activate"`   //1激活 0未激活
	ActivateTime     time.Time `json:"activate_time"` //激活时间
	CreateTime       time.Time `json:"create_time"`
	CodeBackground   string    `json:"code_background"`
	CodeColor        string    `json:"code_color"`
	TimeOut          time.Time `json:"time_out"`
	Status           string    `json:"status"` // success正常 no_use 冻结
	Domain           string    `json:"domain"`
}
