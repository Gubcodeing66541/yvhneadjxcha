package Response

type WeChatAuth struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string
	Openid       string
	Scope        string
	Unionid      string `json:"unionid"`
}

type WeChatAccess struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WeChatUserInfo struct {
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	HeadImgUrl string `headimgurl`
}

type WeChatAuthResp struct {
	Id         int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name       string `json:"name"`
	AppId      string `json:"app_id"`
	AppSecret  string `json:"app_secret"`
	FileName   string `json:"file_name"`
	FileValue  string `json:"file_value"`
	Type       string `json:"type" gorm:"comment:'auth push 对应授权 推送'"` // auth Push
	Url        string `json:"url"`
	UrlSpare   string `json:"url_spare"`
	Status     string `json:"status"` // enable un_enable
	MessageId  string `json:"message_id"`
	CreateTime string `json:"create_time"`
}
