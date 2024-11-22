package Service

import (
	"time"
)

type ServiceRoom struct {
	Id             int    `gorm:"primary_key;AUTO_INCREMENT"`
	RoomId         string `gorm:"index:room_id_idx"`
	ServiceId      int    `gorm:"index:service_id_idx"`
	UserId         int    `gorm:"index:user_id_idx"`
	LateId         int
	LateType       string
	LateMsg        string `json:"late_msg"`
	LateRole       string
	LateIp         string `json:"late_ip"`
	LateUserReadId int
	UserNoRead     int
	ServiceNoRead  int
	IsTop          int `gorm:"comment:'是否置顶'"`
	IsBlack        int `gorm:"comment:'是否拉黑'"`
	IsDelete       int `gorm:"comment:'是否删除'"`
	CreateTime     time.Time
	UpdateTime     time.Time
	Rename         string
}
