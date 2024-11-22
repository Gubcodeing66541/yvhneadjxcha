package Common

type Rename struct {
	Id     int    `json:"id"  gorm:"primary_key;AUTO_INCREMENT"`
	Rename string `json:"rename"`
}
