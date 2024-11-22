package Response

type UserLoginLog struct {
	Id         int    `json:"id"`
	UserId     int    `json:"user_id"`
	ServiceId  int    `json:"service_id"`
	Ip         string `json:"ip"`
	Addr       string `json:"addr"`
	CreateTime string `json:"create_time"`
}

type UserWebAuth struct {
	Code string `json:"code"  `
}
