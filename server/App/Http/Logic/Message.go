package Logic

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"server/App/Common"
	"server/App/Http/Tel"
	Message2 "server/App/Model/Message"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Message struct{}

func (Message) SendToUser(serviceId int, UserId int, Type string, Content string, isPushServiceWs bool) error {
	RoomId := Common.Tools{}.ConvertUserMessageRoomId(serviceId, UserId)
	userIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "user", UserId))
	model := &Message2.Message{
		RoomId:     RoomId,
		From:       serviceId,
		To:         UserId,
		Type:       Type,
		Content:    Content,
		SendRole:   "service",
		CreateTime: time.Now(),
		IsRead:     userIsOnline,
		UserId:     UserId,
		ServiceId:  serviceId,
	}
	Base.MysqlConn.Create(&model)
	if model.Id == 0 {
		return errors.New("消息发送失败")
	}

	update := gin.H{
		"late_msg": Content, "update_time": time.Now(), "late_id": model.Id, "user_no_read": 0,
		"service_no_read": 0, "late_type": Type, "late_role": "service", "is_delete": 0,
	}
	if userIsOnline == 0 {
		update["user_no_read"] = gorm.Expr("user_no_read + ?", 1)
	} else {
		update["late_user_read_id"] = model.Id
	}
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id  = ? and user_id = ?", serviceId, UserId).Updates(update)

	sendMsg := Tel.SocketMessage{
		Id:     model.Id,
		RoomId: RoomId, From: UserId, To: serviceId, Type: Type, Content: Content, ServiceId: serviceId,
		SendRole: "service", CreateTime: time.Now().Format("2006-01-02 15:04:05"), IsRead: userIsOnline, UserId: UserId,
	}

	if isPushServiceWs {
		Common.ApiResponse{}.SendMsgToService(serviceId, "message", sendMsg)
	}
	Common.ApiResponse{}.SendMsgToUser(UserId, "message", sendMsg)

	// 返回OK信息
	return nil
}

func (Message) BotSendToUser(serviceId int, UserId int, Type string, Content string, isPushServiceWs bool) error {
	RoomId := Common.Tools{}.ConvertUserMessageRoomId(serviceId, UserId)
	userIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "user", UserId))
	model := &Message2.Message{
		RoomId:     RoomId,
		From:       serviceId,
		To:         UserId,
		Type:       Type,
		Content:    Content,
		SendRole:   "bot",
		CreateTime: time.Now(),
		IsRead:     userIsOnline,
		UserId:     UserId,
		ServiceId:  serviceId,
	}
	Base.MysqlConn.Create(&model)
	if model.Id == 0 {
		return errors.New("消息发送失败")
	}

	//更新最后一条信息
	update := gin.H{
		"late_msg": Content, "update_time": time.Now(), "late_id": model.Id, "user_no_read": 0,
		"service_no_read": 0, "late_type": Type, "late_role": "service", "is_delete": 0,
	}
	if userIsOnline == 0 {
		update["user_no_read"] = gorm.Expr("user_no_read + ?", 1)
	} else {
		update["late_user_read_id"] = model.Id
	}
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", serviceId, UserId).Updates(update)

	sendMsg := Tel.SocketMessage{
		Id:     model.Id,
		RoomId: RoomId, From: UserId, To: serviceId, Type: Type, Content: Content, ServiceId: serviceId,
		SendRole: "bot", CreateTime: time.Now().Format("2006-01-02 15:04:05"), IsRead: userIsOnline, UserId: UserId,
	}

	if isPushServiceWs {
		Common.ApiResponse{}.SendMsgToService(serviceId, "message", sendMsg)
	}
	Common.ApiResponse{}.SendMsgToUser(UserId, "message", sendMsg)

	// 返回OK信息
	return nil
}
