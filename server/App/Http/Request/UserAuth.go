package Request

type WeChatAuth struct {
	State string `json:"state" uri:"state" form:"state" binding:"required"`
	Code   string `json:"code" uri:"code"  form:"code" binding:"required"`
}

type Token struct {
	Token string `json:"token" uri:"token" form:"token" binding:"required"`
}
