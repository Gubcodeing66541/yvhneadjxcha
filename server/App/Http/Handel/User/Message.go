package User

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	Message2 "server/App/Model/Message"
	"server/Base"
	"time"
)

type Message struct{}

func (Message) SendToActivate(c *gin.Context) {
	var req Request.UserSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" {
		Common.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	RoleId := Common.Tools{}.GetRoleId(c)
	ActivateId := Common.Tools{}.GetServiceId(c)
	RoomId := Common.Tools{}.ConvertUserMessageRoomId(ActivateId, RoleId)
	model := &Message2.Message{
		RoomId: RoomId, From: RoleId, To: ActivateId, Type: req.Type, Content: req.Content,
		SendRole: "service", CreateTime: time.Now(), IsRead: 0, UserId: RoleId}
	Base.MysqlConn.Create(&model)
	if model.Id == 0 {
		Common.ApiResponse{}.Error(c, "消息发送失败", gin.H{})
		return
	}

	res, err := json.Marshal(model)
	if err != nil {
		Common.ApiResponse{}.Error(c, "消息发送失败.", gin.H{})
		return
	}

	// 给客服
	serviceIdName := Common.Tools{}.GetWebSocketId(c)
	Base.WebsocketHub.SendToUserId(serviceIdName, res)

	// 给用户推送
	userIdName := fmt.Sprintf("%s:%d", "user", RoleId)
	Base.WebsocketHub.SendToUserId(userIdName, res)

	Common.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})
}

func (Message) List(c *gin.Context) {
	var req Request.MsgUserId
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	var message []Message2.Message
	Base.MysqlConn.Find(&message, "service_id = ? and user_id = ?", roleId, req.UserId)
	Common.ApiResponse{}.Success(c, "用户不存在", gin.H{"list": message})
}
