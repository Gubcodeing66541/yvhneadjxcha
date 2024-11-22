package Service

import "time"

type Order struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int       `json:"service_id"`
	Amount     int       `json:"amount"`
	Day        int       `json:"day"`
	CreateTime time.Time `json:"createTime"`
}
