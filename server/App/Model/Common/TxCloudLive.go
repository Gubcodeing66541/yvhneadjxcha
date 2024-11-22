package Common

type TxCloudLive struct {
	Id      int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Domain  string `json:"domain"`
	AppName string `json:"domain"`
}
