package ServiceManager

import (
	"time"
)

type ServiceManagerAuth struct {
	ServiceManagerId int `json:"service_manager_id"`
	Mobile           string
	Username         string
	Password         string
	TimeOut          time.Time
	CreateTime       time.Time
	UpdateTime       time.Time
}
