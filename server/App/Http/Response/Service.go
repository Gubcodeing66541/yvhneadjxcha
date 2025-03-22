package Response

import (
	"server/App/Constant"
	"time"
)

type ServiceInfo struct {
	ServiceId      int       `json:"service_id"`
	Name           string    `json:"name"`
	Head           string    `json:"head"`
	Mobile         string    `json:"mobile"`
	Username       string    `json:"username"`
	Type           string    `json:"type"` //auth push
	Code           string    `json:"code"`
	Domain         string    `json:"domain"`
	Web            string    `json:"web"`
	Host           string    `json:"host"`
	TimeOut        string    `json:"time_out"`
	CreateTime     time.Time `json:"create_time"`
	CodeBackground string    `json:"code_background"`
	CodeColor      string    `json:"code_color"`
	RoomCount      int       `json:"room_count"`
	BlackCount     int       `json:"black_count"`
	BotHead        string    `json:"bot_head"`
	BindDomain     string    `json:"bind_domain"`
	BindAction     string    `json:"bind_action"`
}

type RespServiceList struct {
	ServiceId int                 `json:"service_id"`
	Name      string              `json:"name"`
	Head      string              `json:"head"`
	Username  string              `json:"username"`
	Type      string              `json:"type"`
	Role      string              `json:"role"`
	Code      string              `json:"code"`
	TimeOut   Constant.SystemTime `json:"time_out"`
	Domain    string              `json:"domain"`
	Status    string              `json:"status"`
}

type ServiceUserCount struct {
	AllCnt  int    `json:"all_cnt"`
	UserCnt int    `json:"user_cnt"`
	Dates   string `json:"dates"`
	Rate    int    `json:"rate"`
}

type ServiceUserCountStatistics struct {
	AllCnt  int    `json:"all_count"`
	UserCnt int    `json:"service_count"`
	Dates   string `json:"dates"`
	Rate    int    `json:"rate"`
}

// 客服表
type ServiceList struct {
	Id               int                 `json:"id" gorm:"primary_key;AUTO_INCREMENT" `
	ServiceId        int                 `json:"service_id"`
	ServiceManagerId int                 `json:"service_manager_id"`
	Name             string              `json:"name"`
	Head             string              `json:"head"`
	Mobile           string              `json:"mobile"`
	Username         string              `json:"username"`
	Type             string              `json:"type"` //auth push
	Role             string              `json:"role"` //user  vip
	Day              int                 `json:"day"`
	Code             string              `json:"code"`
	IsActivate       int                 `json:"is_activate"`   //1激活 0未激活
	ActivateTime     Constant.SystemTime `json:"activate_time"` //激活时间
	CreateTime       Constant.SystemTime `json:"create_time"`
	CodeBackground   string              `json:"code_background"`
	CodeColor        string              `json:"code_color"`
	TimeOut          Constant.SystemTime `json:"time_out"`
	UserCnt          int                 `json:"user_cnt"`
	IsOnline         int                 `json:"is_online"` // 1在线 0不在线
	Status           string              `json:"status"`    //no_active未激活 time_out过期 success 正常 no_use冻结
	Domain           string              `json:"domain"`    // 域名
}

// 客服表
type ServiceListBySocket struct {
	Id               int       `json:"id" gorm:"primary_key;AUTO_INCREMENT" `
	ServiceId        int       `json:"service_id"`
	ServiceManagerId int       `json:"service_manager_id"`
	Name             string    `json:"name"`
	Head             string    `json:"head"`
	Mobile           string    `json:"mobile"`
	Username         string    `json:"username"`
	Type             string    `json:"type"` //auth push
	Role             string    `json:"role"`
	Day              int       `json:"day"`
	Code             string    `json:"code"`
	IsActivate       int       `json:"is_activate"`   //1激活 0未激活
	ActivateTime     time.Time `json:"activate_time"` //激活时间
	CreateTime       time.Time `json:"create_time"`
	CodeBackground   string    `json:"code_background"`
	CodeColor        string    `json:"code_color"`
	TimeOut          time.Time `json:"time_out"`
	Status           string    `json:"status"` // success正常 no_use 冻结
	Member           string    `json:"member"`
}
