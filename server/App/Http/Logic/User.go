package Logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	Service2 "server/App/Model/Service"
	ServiceManager2 "server/App/Model/ServiceManager"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type User struct{}

func (User) WeChatAuth(req Request.WeChatAuth) error {
	// 微信code换的信息
	return nil
}

func (User) OpenIdToUser(openId string) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "open_id = ?", openId)
	return user
}

func (User) UnionIdToUser(unionId string) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "union_id = ?", unionId)
	return user
}

func (User) UserIdToUser(userId int) User2.User {
	var user User2.User
	Base.MysqlConn.Find(&user, "user_id = ?", userId)
	return user
}

// CookieUUIDToUser 通过系统生成的uuid获取用户信息
func (User) CookieUUIDToUser(cookieUuid string) User2.User {
	var user User2.User
	var userMap User2.UserAuthMap
	Base.MysqlConn.Find(&userMap, "cookie_uid = ?", cookieUuid)
	Base.MysqlConn.Find(&user, "user_Id = ?", userMap.UserId)
	return user
}

func (User) CreateUser(openId string, Name string, Header string, Sex int, unionId string) User2.User {
	now := time.Now()
	user := User2.User{OpenId: openId, UserName: Name, UserHead: Header, CreateTime: now, UpdateTime: now, UnionId: unionId}
	Base.MysqlConn.Create(&user)
	return user
}

func (User) CreateWebUser(Name string, Header string, Token string) User2.User {
	now := time.Now()
	user := User2.User{UserName: Name, UserHead: Header, CreateTime: now, UpdateTime: now, Token: Token}
	Base.MysqlConn.Create(&user)
	return user
}

func (User) HandelLeaveMessage(c *gin.Context, serviceId int, userId int) {
	// 如果客服在线则不用管了
	ServiceIsOnline := Base.WebsocketHub.UserIdIsOnline(Common.Tools{}.GetServiceWebSocketId(serviceId))
	if ServiceIsOnline == 1 {
		return
	}

	var leaveMsg []Service2.ServiceMessage
	Base.MysqlConn.Find(&leaveMsg, "service_id = ? and type = 'leave'", serviceId)

	for key, v := range leaveMsg {
		if v.Status == "enable" {
			_ = Message{}.SendToUser(serviceId, userId, leaveMsg[key].MsgType, leaveMsg[key].MsgInfo, true)
		}
	}
}

// 处理机器人消息
func (User) HandelBotMessage(c *gin.Context, serviceId int, userId int, req Request.UserSendMessage) {
	service, _ := Service{}.IdGet(serviceId)
	if service.ServiceId == 0 {
		return
	}

	// 检测是否启动bot
	var bot ServiceManager2.ServiceManagerBot
	Base.MysqlConn.Find(&bot, "service_manager_id = ?", service.ServiceManagerId)
	if bot.Status != "run" {
		return
	}

	// 检测是否出发机器人消息
	var Msg ServiceManager2.ServiceManagerBotMessage
	Base.MysqlConn.Find(&Msg, "service_manager_id = ? && problem like ?", service.ServiceManagerId, "%"+req.Content+"%")
	if Msg.Id != 0 {
		err := Message{}.BotSendToUser(serviceId, userId, "text", Msg.Answer, true)
		if err != nil {
			fmt.Println("bot 推送 err", serviceId, userId, "text", Msg.Answer, err.Error())
		}
	}
}
