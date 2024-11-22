package ServiceManager

import "time"

type ServiceManager struct {
	ServiceManagerId  int       `gorm:"primary_key;AUTO_INCREMENT" json:"service_manager_id"`
	Member            string    `json:"member"`
	Password          string    `json:"password"`
	Name              string    `json:"name"`
	Head              string    `json:"head"`
	CreateTime        time.Time `json:"create_time"`
	UpdateTime        time.Time `json:"update_time"`
	Ip                string    `json:"ip"`
	IdCardA           string    `json:"id_card_a"`
	IdCardB           string    `json:"id_card_b"`
	RealName          string    `json:"real_name"`
	IdCardNumber      string    `json:"id_card_number"`
	IdCardApplyStatus string    `json:"id_card_apply_status"` // 表示 none没有上传  apply 审批中 agree 同意 refuse 拒绝
	Account           int       `json:"account"`              // 钱包金额最小到分数 *100是元
	Status            string    `json:"status"`               // success正常 no_use 冻结

}
