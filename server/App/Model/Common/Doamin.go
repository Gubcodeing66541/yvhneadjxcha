package Common

import "time"

type Domain struct {
	Id              int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Domain          string    `json:"domain" gorm:"type:varchar(500);"`
	Type            string    `json:"type"`               // public private action
	BindServiceId   int       `json:"bind_service_id"`    // bind service id
	WeChatBanStatus string    `json:"we_chat_ban_status"` // success error ban
	Status          string    `json:"status"`             // enable un_enable down下架
	BindCnt         int       `json:"bind_cnt"`           // 绑定数量
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}
