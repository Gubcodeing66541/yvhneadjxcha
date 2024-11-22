package Count

import "time"

type CountServiceRoom struct {
	Id               int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ServiceManagerId int       `json:"service_manager_id"`
	ServiceId        int       `json:"service_id"` //0是汇总统计 非0是客服统计
	AddUser          int       `json:"add_user"`
	AllUser          int       `json:"all_user"`
	OnlineUser       int       `json:"online_user"`
	OnlineService    int       `json:"online_service"`
	CountTime        time.Time `json:"count_time"`
	CreateTime       time.Time `json:"create_time"`
}
