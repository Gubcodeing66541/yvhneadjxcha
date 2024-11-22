package Common

type Token struct {
	Id    int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Token string `json:"token"`
}
