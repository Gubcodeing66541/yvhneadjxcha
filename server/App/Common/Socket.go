package Common

import (
	"server/Base"
)

type Socket struct{}

func (Socket) GetAll() (service []int, user []int) {
	return Base.WebsocketHub.GetAllConn()
}
