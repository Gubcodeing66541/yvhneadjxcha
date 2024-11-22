package Service

import "time"

type ServiceAuth struct {
	ServiceId        int `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceManagerId int `json:"service_manager_id"`
	Mobile           string
	Username         string
	Password         string
	TimeOut          time.Time
	CreateTime       time.Time
	UpdateTime       time.Time
}
