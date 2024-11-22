package Service

import "time"

type ServiceRoomDetail struct {
	Id         int    `gorm:"primary_key;AUTO_INCREMENT"`
	RoomId     string `gorm:"index:room_id_idx"`
	UserId     int    `gorm:"index:user_id_idx"`
	ServiceId  int    `gorm:"index:service_id_idx"`
	Drive      string
	IP         string
	Map        string
	Mobile     string
	Tag        string
	CreateTime time.Time
}
