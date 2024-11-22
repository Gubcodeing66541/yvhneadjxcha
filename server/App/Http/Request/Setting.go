package Request

import "time"

type SettingPageLimit struct {
	Search string `json:"search" uri:"search" xml:"search" form:"search"`
	Type   string `json:"type" uri:"type" xml:"type" form:"type"`
	Page   int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type CreateSetting struct {
	Type       string    `json:"type" uri:"type" xml:"type" form:"type" binding:"required"`
	Value      string    `json:"value" uri:"value" xml:"value" form:"value" binding:"required"`
	CreateTime time.Time `json:"create_time" uri:"create_time" xml:"create_time" form:"create_time"`
}

type CreateSettingUpdate struct {
	Id         int       `json:"id" uri:"id" xml:"id" form:"id" binding:"required"`
	Type       string    `json:"type" uri:"type" xml:"type" form:"type" binding:"required"`
	Value      string    `json:"value" uri:"value" xml:"value" form:"value" binding:"required"`
	CreateTime time.Time `json:"create_time" uri:"create_time" xml:"create_time" form:"create_time"`
}
