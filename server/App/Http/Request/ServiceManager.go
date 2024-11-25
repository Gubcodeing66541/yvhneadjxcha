package Request

import "time"

type ServiceManagerId struct {
	ServiceManagerId int    `json:"service_manager_id" uri:"service_manager_id" form:"service_manager_id" binding:"required"`
	StartTime        string `json:"start_time" uri:"start_time" xml:"start_time" form:"start_time"`
	EndTime          string `json:"end_time" uri:"end_time" xml:"end_time" form:"end_time"`
	Search           string `json:"search" uri:"search" xml:"search" form:"search"`
	Page             int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset           int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
	Status           string `json:"status" uri:"status" xml:"status" form:"status"`
}

type ServiceManagerUpdate struct {
	Name              string `json:"name" uri:"name" form:"name"`
	Head              string `json:"head" uri:"head" form:"head"`
	IdCardA           string `json:"id_card_a" uri:"id_card_a" form:"id_card_a"`
	IdCardB           string `json:"id_card_b" uri:"id_card_b" form:"id_card_b"`
	RealName          string `json:"real_name" uri:"real_name" form:"real_name"`
	IdCardNumber      string `json:"id_card_number" uri:"id_card_number" form:"id_card_number"`
	IdCardApplyStatus string `json:"id_card_apply_status" uri:"id_card_apply_status" form:"id_card_apply_status"`
}

type ServiceManagerMessageAddOrUpdate struct {
	Id               int       `json:"id" uri:"id" form:"id"`
	Title            string    `json:"title" uri:"title" form:"title"`
	Content          string    `json:"content" uri:"content" form:"content" binding:"required"`
	Type             string    `json:"type" uri:"type" form:"type" binding:"required"`
	AddServiceId     int       `json:"add_service_id"`
	Status           string    `json:"status"` // 表示  apply 审批中 agree 同意 refuse 拒绝
	UpdateTime       time.Time `json:"update_time" `
	ServiceManagerId int       `json:"service_manager_id"`
}

type ServiceManagerMessageDelete struct {
	Id int `json:"id" uri:"id" form:"id" binding:"required"`
}

type MemberCreateService struct {
	ServiceName string `json:"service_name" uri:"service_name" form:"service_name" binding:"required"`
	Day         int    `json:"day" uri:"day" form:"day"  binding:"required"`
}

type MemberCreateServiceList struct {
	ServiceNumber int `json:"service_number" uri:"service_number" form:"service_number"   binding:"required"`
	Day           int `json:"day" uri:"day" form:"day"  binding:"required"`
}

type MemberCreateServiceUpdate struct {
	ServiceId int    `json:"service_id" uri:"service_id" form:"service_id"   binding:"required"`
	Name      string `json:"name" uri:"name" form:"name"   binding:"required"`
}

type MemberCreateServiceId struct {
	ServiceId int `json:"service_id" uri:"service_id" form:"service_id"   binding:"required"`
}

type MemberServiceRenewal struct {
	Username string `json:"username" uri:"username" form:"username" binding:"required"`
	Day      uint   `json:"day" uri:"day" form:"day"  binding:"required"`
}

type MemberServiceRenewalList struct {
	UsernameList string `json:"username_list" uri:"username_list" form:"username_list" binding:"required"`
	Day          uint   `json:"day" uri:"day" form:"day"  binding:"required"`
}

type BlackPage struct {
	StartTime   string `json:"start_time"  uri:"start_time" xml:"start_time" form:"start_time"`
	EndTime     string `json:"end_time"  uri:"end_time" xml:"end_time" form:"end_time"`
	UserName    string `json:"user_name"  uri:"user_name" xml:"user_name" form:"user_name"`
	ServiceName string `json:"service_name"  uri:"service_name" xml:"service_name" form:"service_name"`
	Search      string `json:"search" uri:"search" xml:"search" form:"search"`
	Page        int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset      int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type ServiceManagerReset struct {
	ServiceManagerId int    `json:"service_manager_id" uri:"service_manager_id" form:"service_manager_id" binding:"required"`
	Name             string `json:"name" uri:"name" form:"name" binding:"required"`
	Password         string `json:"password" uri:"password" form:"password"`
}
