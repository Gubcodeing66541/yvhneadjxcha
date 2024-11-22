package Service

type ServiceNoticeSetting struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ServiceId int    `json:"service_id"`
	IsShow    int    `json:"is_show"`
	Image     string `json:"image"`
	Text      string `json:"text"`
}
