package Config

type Oss struct {
	Qiniu Qiniu `json:"qiniu"`
}

type Qiniu struct {
	AccessKey string	`json:"access_key"`
	SecretKey string	`json:"secret_key"`
	Bucket    string	`json:"bucket"`
}
