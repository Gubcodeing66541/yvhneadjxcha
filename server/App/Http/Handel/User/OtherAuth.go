package User

import (
	"fmt"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	Message2 "server/App/Model/Message"
	"server/App/Model/Service"
	"server/App/Model/User"
	"server/Base"
	"time"

	"github.com/gin-gonic/gin"
)

type otherAuth struct{}

var OtherAuth = otherAuth{}

// 首先是同一个域名 有token就有，没有就注册，就这么简单 仅仅只有入口请求落地是可以直接拿到的才对
func (otherAuth) Action(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	service, err := Logic.Service{}.Get(req.Code)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.Code+"客服不存在", gin.H{})
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 准备绑定的用户
	var userModel User.User

	// 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	username := fmt.Sprintf("%s", Common.Tools{}.GetRename())
	userModel = Logic.User{}.CreateUser("", username, Common.Tools{}.GetHead(), 0, "")

	_ = Logic.ServiceRoom{}.Get(userModel, service.ServiceId, c.ClientIP(), Common.ClientAgentTools{}.GetDrive(c))
	token := Common.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0)

	ip := c.ClientIP()
	var black Service.ServiceBlack
	Base.MysqlConn.Find(
		&black,
		"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?)",
		service.ServiceId, ip, service.ServiceId, userModel.UserId)
	if black.Id != 0 {
		Common.ApiResponse{}.Error(c, "无法访问", gin.H{})
		return
	}

	// 上线更新update时间和未读
	update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
	Base.MysqlConn.Model(&Service.ServiceRoom{}).
		Where("user_id = ? and service_id = ? ", userModel.UserId, service.ServiceId).
		Updates(update)

	Base.MysqlConn.Model(&Request.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"ip": c.ClientIP()})

	// 所有消息已读
	Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"is_read": 1})

	Base.MysqlConn.Model(&Request.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"ip": c.ClientIP()})

	Base.MysqlConn.Create(&User.UserLoginLog{
		UserId:     userModel.UserId,
		ServiceId:  service.ServiceId,
		Ip:         c.ClientIP(),
		Addr:       "",
		CreateTime: time.Now(),
	})

	action := Logic.Domain{}.GetAction()
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token, "action": action})
}
