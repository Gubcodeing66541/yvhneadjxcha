package Request

type WeChatList struct {
	Type string `json:"type"  form:"type"  uri:"type" ` // auth Push
	Page Page
}

type WeChatType struct {
	Search string `json:"search" uri:"search" xml:"search" form:"search"`
	Page   int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
	Type   string `json:"type"  form:"type"  uri:"type"` // auth Push
}

type OpenOrClose struct {
	Id     int    `json:"id" form:"id"  uri:"id" binding:"required"`
	Status string `json:"status" form:"status"  uri:"status"  binding:"required"` // enable un_enable
}

type CreateWeChatAuth struct {
	Id        int    `json:"id" form:"id" json:"id" uri:"id" `
	Name      string `json:"name" form:"name" uri:"name" binding:"required"`
	AppId     string `json:"app_id" form:"app_id" uri:"app_id" binding:"required"`
	AppSecret string `json:"app_secret" form:"app_secret"  uri:"app_secret" binding:"required"`
	FileName  string `json:"file_name" form:"file_name"  uri:"file_name" binding:"required"`
	FileValue string `json:"file_value" form:"file_value"  uri:"file_value"  binding:"required"`
	Type      string `json:"type"  form:"type"  uri:"type" binding:"required"` // auth Push
	Url       string `json:"url" form:"url"  uri:"url" xml:"url" binding:"required"`
	UrlSpare  string `json:"url_spare" form:"url_spare"`
	MessageId string `json:"message_id" form:"message_id"`
	Status    string `json:"status" form:"status"` // enable un_enable
}

type UpdateWeChatAuth struct {
	Id        int    `json:"id" form:"id" json:"id" uri:"id" binding:"required"`
	Name      string `json:"name" form:"name" uri:"name"`
	AppId     string `json:"app_id" form:"app_id" uri:"app_id" `
	AppSecret string `json:"app_secret" form:"app_secret" uri:"app_secret" `
	FileName  string `json:"file_name" form:"file_name" uri:"file_name" `
	FileValue string `json:"file_value" form:"file_value" uri:"file_value"  `
	Type      string `json:"type"  form:"type"  uri:"type" ` // auth Push
	Url       string `json:"url" form:"url"  uri:"url" xml:"url" `
	UrlSpare  string `json:"url_spare" form:"url_spare" `
	MessageId string `json:"message_id" form:"message_id"  uri:"message_id" xml:"message_id" `
	Status    string `json:"status" form:"status" ` // enable un_enable
}

type DeleteWeChatAuth struct {
	Id     int    `json:"id" `
	Status string `json:"status"` // enable un_enable
}

type WeChatSwitch struct {
	Id       int    `json:"id" form:"id" json:"id" uri:"id" binding:"required"`
	Url      string `json:"url" form:"url"  uri:"url" xml:"url"  binding:"required"`
	UrlSpare string `json:"url_spare" form:"url_spare"  binding:"required"`
}
