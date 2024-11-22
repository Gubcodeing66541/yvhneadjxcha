package ServiceManager

import (
	"time"
)

type ServiceManagerRenewRecorder struct {
	Id                   int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	OrderId              string    `json:"order_id"`
	ServiceManagerId     int       `json:"service_manager_id"`
	ServiceManagerMember string    `json:"service_manager_member"`
	Member               string    `json:"member"`
	ServiceId            int       `json:"service_id"`
	OldAccount           int       `json:"old_account"`
	Account              int       `json:"account"`
	Renew                int       `json:"renew"`
	Reason               string    `json:"reason"`
	PayType              string    `json:"pay_type"`
	CreateTime           time.Time `json:"create_time"`
}
