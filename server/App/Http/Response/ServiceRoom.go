package Response

import (
	"server/App/Constant"
)

type ServiceRoom struct {
	Id            int                 `json:"id"`
	LateType      string              `json:"late_type"`
	UserNoRead    int                 `json:"user_no_read"`
	ServiceNoRead int                 `json:"service_no_read"`
	LateMsg       string              `json:"late_msg"`
	IsOnline      int                 `json:"is_online"`
	IsTop         int                 `json:"is_top"`
	UpdateTime    Constant.SystemTime `json:"update_time"`
	UserId        int                 `json:"user_id"`
	UserName      string              `json:"user_name"`
	UserHead      string              `json:"user_head"`
	Rename        string              `json:"rename"`
}

type ServiceRoomListRes struct {
	Id            int    `json:"id"`
	LateType      string `json:"late_type"`
	UserNoRead    int    `json:"user_no_read"`
	ServiceNoRead int    `json:"service_no_read"`
	LateMsg       LateMsg
	IsOnline      int                 `json:"is_online"`
	IsTop         int                 `json:"is_top"`
	UpdateTime    Constant.SystemTime `json:"update_time"`
	UserId        int                 `json:"user_id"`
	UserName      string              `json:"user_name"`
	UserHead      string              `json:"user_head"`
	Rename        string              `json:"rename"`
}

type LateMsg struct {
	Text string `json:"text"`
	Img  string `json:"imt"`
}

type UserDetail struct {
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	UserHead string `json:"user_head"`
	Rename   string `json:"rename"`
	IsTop    string `json:"is_top"`
	Drive    string `json:"drive"`
	Ip       string `json:"ip"`
	Map      string `json:"map"`
	Mobile   string `json:"mobile"`
	Wechat   string `json:"wechat"`
	Age      string `json:"age"`
	Tag      string `json:"tag"`
}

type UserBlackList struct {
	Id         int                 `gorm:"primary_key;AUTO_INCREMENT"`
	ServiceId  int                 `gorm:"index:service_id_idx" json:"service_id"`
	Type       string              `json:"type"` //拉黑类型ip user
	UserId     int                 `json:"user_id"`
	Ip         string              `json:"ip"`
	UserName   string              `json:"user_name"`
	UserHead   string              `json:"user_head"`
	CreateTime Constant.SystemTime `json:"create_time"`
}

type DeleteUserDay struct {
	Day int `json:"day" uri:"day" form:"day" `
}

type ServiceRoomList struct {
	UserId      int                 `json:"user_id"`
	UserName    string              `json:"user_name"`
	UserHead    string              `json:"user_head"`
	Ip          string              `json:"ip"`
	Map         string              `json:"map"`
	Drive       string              `json:"drive"`
	Mobile      string              `json:"mobile"`
	Wechat      string              `json:"wechat"`
	Tag         string              `json:"tag"`
	Name        string              `json:"name"`
	IsOnline    int                 `json:"is_online"`
	ServiceId   int                 `json:"service_id"`
	ServiceHead string              `json:"service_head"`
	CreateTime  Constant.SystemTime `json:"create_time"`
}

type IpCount struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
