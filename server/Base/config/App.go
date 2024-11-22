package Config

type App struct {
	Debug           bool `json:"debug"`
	LinkIsShowImage bool `json:"link_is_show_image"`
	Database        Database
	Client          Client
	Oss             Oss
	Mq              Mq
	HeadImgUrl      string  `json:"head_img_url"`
	Manager         Manager `json:"manager"`
	HttpHost        string  `json:"http_host"`
	Check           check   `json:"check"`
}

type check struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}
