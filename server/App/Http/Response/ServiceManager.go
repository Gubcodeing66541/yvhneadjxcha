package Response

import (
	"server/App/Constant"
)

type ServiceManagerInfo struct {
	ServiceManagerId  int                 `gorm:"primary_key;AUTO_INCREMENT" json:"service_manager_id"`
	Name              string              `json:"name"`
	Head              string              `json:"head"`
	CreateTime        Constant.SystemTime `json:"create_time"`
	UpdateTime        Constant.SystemTime `json:"update_time"`
	Ip                string              `json:"ip"`
	IdCardA           string              `json:"id_card_a"`
	IdCardB           string              `json:"id_card_b"`
	RealName          string              `json:"real_name"`
	IdCardNumber      string              `json:"id_card_number"`
	IdCardApplyStatus string              `json:"id_card_apply_status"` // 表示 none没有上传  apply 审批中 agree 同意 refuse 拒绝
	Account           int                 `json:"account"`              // 钱包金额最小到分数 *100是元
	Status            string              `json:"status"`
}

type ServiceManagerDetail struct {
	ServiceManagerId  int                 `gorm:"primary_key;AUTO_INCREMENT" json:"service_manager_id"`
	Member            string              `json:"member"`
	Password          string              `json:"password"`
	Name              string              `json:"name"`
	Head              string              `json:"head"`
	CreateTime        Constant.SystemTime `json:"create_time"`
	UpdateTime        Constant.SystemTime `json:"update_time"`
	Ip                string              `json:"ip"`
	IdCardA           string              `json:"id_card_a"`
	IdCardB           string              `json:"id_card_b"`
	RealName          string              `json:"real_name"`
	IdCardNumber      string              `json:"id_card_number"`
	IdCardApplyStatus string              `json:"id_card_apply_status"` // 表示 none没有上传  apply 审批中 agree 同意 refuse 拒绝
	Account           int                 `json:"account"`              // 钱包金额最小到分数 *100是元
	ServiceCnt        int                 `json:"service_cnt"`          // 客服数量
	Status            string              `json:"status"`
}

type ServiceBlack struct {
	Id         int                 `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int                 `gorm:"index:service_id_idx" json:"service_id"`
	Type       string              `json:"type"` //拉黑类型ip user
	UserId     int                 `json:"user_id"`
	Ip         string              `json:"ip"`
	UserName   string              `json:"user_name"`
	UserHead   string              `json:"user_head"`
	Name       string              `json:"name"`
	Username   string              `json:"username"`
	CreateTime Constant.SystemTime `json:"create_time"`
}

type ServiceManagerCount struct {
	RoomCnt    int `json:"room_cnt"`
	MessageCnt int `json:"message_cnt"`
}

type ServiceManagerRoomCount struct {
	ServiceId int    `json:"service_id"`
	Name      string `json:"name"`
	RoomCnt   int    `json:"room_cnt"`
}

type ServiceManagerMessageCount struct {
	ServiceId  int    `json:"service_id"`
	Name       string `json:"name"`
	MessageCnt int    `json:"message_cnt"`
}

type ServiceManagerRenewRecorder struct {
	Id                   int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	OrderId              string `json:"order_id"`
	ServiceManagerId     int    `json:"service_manager_id"`
	ServiceManagerMember string `json:"service_manager_member"`
	Member               string `json:"member"`
	ServiceId            int    `json:"service_id"`
	OldAccount           int    `json:"old_account"`
	Account              int    `json:"account"`
	Renew                int    `json:"renew"`
	Reason               string `json:"reason"`
	PayType              string `json:"pay_type"`
	CreateTime           string `json:"create_time"`
}

type SocketDataResponse struct {
	AllUserCnt      int `json:"all_user_cnt"`
	AllMessageCnt   int `json:"all_message_cnt"`
	TodayUserCnt    int `json:"today_user_cnt"`
	TodayMessageCnt int `json:"today_message_cnt"`
	OnlineUser      int `json:"online_user"`
	OnlineService   int `json:"online_service"`
}

type RegionResponse struct {
	Map string `json:"map"`
	Cnt int    `json:"cnt"`
}

type SocketUserResponse struct {
	UserName    string `json:"user_name"`
	UserHead    string `json:"user_head"`
	Map         string `json:"map"`
	ServiceName string `json:"service_name"`
}

type ServiceRank struct {
	Name string `json:"name"`
	Head string `json:"head"`
	Cnt  int    `json:"cnt"`
}

type ContTime struct {
	Time string `json:"time"`
	Cnt  int    `json:"cnt"`
}

type ServiceManagerBotMessageResp struct {
	Id               int    `json:"id" uri:"id"  form:"id" `
	ServiceManagerId int    `json:"service_manager_id"`
	Problem          string `json:"problem" uri:"problem"  form:"problem"`
	Answer           string `json:"answer" uri:"answer"  form:"answer"`
	CreateTime       string `json:"create_time"`
}
