package Response

import (
	"server/App/Http/Request"
	"time"
)

type ResultList struct {
	Data interface{} `json:"data"`
	Page Request.Page
}

type ResConfig struct {
	SystemName     string
	SystemVideo    string
	SystemVersion  string
	SystemLogo     string
	PcDownloadLink string
	AdDefault      string
	RealName       int
}

type Count struct {
	Counted int
}

type Setting struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Type       string `json:"type"` // keyword 关键字 notice  公告
	Value      string `json:"value"`
	CreateTime string `json:"create_time"`
}

type SettingResponse struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Type       string `json:"type"` // keyword 关键字 notice  公告
	Value      string `json:"value"`
	CreateTime string `json:"create_time"`
}

type ServiceManagerMessage struct {
	Id               int       `json:"id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Type             string    `json:"type"`
	ServiceManagerId int       `json:"service_manager_id"`
	AddServiceId     int       `json:"add_service_id"`
	Status           string    `json:"status"` // 表示  apply 审批中 agree 同意 refuse 拒绝
	CreateTime       time.Time `json:"create_time" `
	UpdateTime       time.Time `json:"update_time"`
}

type ServiceManagerMessageTimeToString struct {
	Id               int    `json:"id"`
	Title            string `json:"title"`
	Content          string `json:"content"`
	Type             string `json:"type"`
	ServiceManagerId int    `json:"service_manager_id"`
	AddServiceId     int    `json:"add_service_id"`
	Status           string `json:"status"` // 表示  apply 审批中 agree 同意 refuse 拒绝
	CreateTime       string `json:"create_time" `
	UpdateTime       string `json:"update_time"`
}

type Counted struct {
	Cnt int `json:"cnt"`
}
