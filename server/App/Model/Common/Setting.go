package Common

import (
	"time"
)

type Setting struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Type       string    `json:"type"` // keyword 关键字 notice  公告
	Value      string    `json:"value"`
	CreateTime time.Time `json:"create_time"`
}
