package Logic

import (
	"server/App/Http/Request"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

// 快捷消息
type ServiceMessage struct {
}

// 创建消息
func (ServiceMessage) Create(serviceId int, req Request.CreateServiceMessage) {
	Base.MysqlConn.Create(&Service2.ServiceMessage{ServiceId: serviceId, Status: "un_enable", MsgInfo: req.MsgInfo, MsgType: req.MsgType, Type: req.Type, CreateTime: time.Now()})
}

// 删除招呼
func (ServiceMessage) Delete(id int, serviceId int) {
	Base.MysqlConn.Delete(&Service2.ServiceMessage{}, "id = ? and service_id = ? ", id, serviceId)
}

// 修改招呼
func (ServiceMessage) Update(id int, msgType string, msgInfo string, serviceId int, status string) {
	var serviceMessage Service2.ServiceMessage
	Base.MysqlConn.Model(&serviceMessage).Where("id = ? and service_id = ?", id, serviceId).Updates(Service2.ServiceMessage{MsgType: msgType, MsgInfo: msgInfo, Status: status})
}

// 招呼列表
func (ServiceMessage) List(serviceId int, typeStr string) []Service2.ServiceMessage {
	var serviceMessage []Service2.ServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ? and type = ?", serviceId, typeStr)
	return serviceMessage
}

func (ServiceMessage) GetById(serviceId int, id int) Service2.ServiceMessage {
	var serviceMessage Service2.ServiceMessage
	Base.MysqlConn.Find(&serviceMessage, "service_id = ? and id = ?", serviceId, id)
	return serviceMessage
}
