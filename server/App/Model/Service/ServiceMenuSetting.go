package Service

type ServiceMenuSetting struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ServiceId int    `json:"service_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Action    string `json:"action"`
	Sort      int    `json:"sort"`
	Tag       string `json:"tag"`
}
