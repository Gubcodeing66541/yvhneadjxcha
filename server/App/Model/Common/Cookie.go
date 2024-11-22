package Common

type Cookie struct {
	Id     string `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserId string `json:"userId"`
	Code   string `json:"code"`
	Cookie string `json:"cookie"`
}
