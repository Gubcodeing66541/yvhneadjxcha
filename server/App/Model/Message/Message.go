package Message

import "time"

type Message struct {
	Id         int       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RoomId     string    `json:"room_id" gorm:"index:room_id_idx"`
	From       int       `json:"from"`
	To         int       `json:"to"`
	Type       string    `json:"type"`
	Content    string    `json:"content" gorm:"type:text"`
	SendRole   string    `json:"send_role"` // user activate
	CreateTime time.Time `json:"create_time"`
	IsRead     int       `json:"is_read"`
	UserId     int       `json:"user_id"`
	ServiceId  int       `json:"service_id"`
	IsDel      int       `json:"is_del"` //1 del
}
