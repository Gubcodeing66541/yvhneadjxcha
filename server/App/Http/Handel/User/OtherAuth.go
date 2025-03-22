package User

import (
	"fmt"
	"net/http"
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
		Uuid string `json:"uuid"`
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

	var uuid = req.Uuid
	if req.Uuid == "" {
		uuid = Common.Tools{}.CreateUserMember()
	}

	Base.MysqlConn.Create(&User.UserLoginLog{
		UserId:     userModel.UserId,
		ServiceId:  service.ServiceId,
		Ip:         c.ClientIP(),
		Addr:       "",
		CreateTime: time.Now(),
	})

	Base.MysqlConn.Create(&User.UserAuthMap{
		UserId: userModel.UserId, CookieUid: uuid,
	})

	action := fmt.Sprintf("%s?uuid=%s&code=%s", Logic.Domain{}.GetAction(), uuid, req.Code)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token, "action": action, "uuid": uuid})
}

func (otherAuth) Domain(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
		Uuid string `json:"uuid"`
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

	roleId := Common.Tools{}.GetRoleId(c)

	// 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	userModel = Logic.User{}.UserIdToUser(roleId)

	_ = Logic.ServiceRoom{}.Get(userModel, service.ServiceId, c.ClientIP(), Common.ClientAgentTools{}.GetDrive(c))

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

	fmt.Println("ok-----------------------------")
	domainInfo := fmt.Sprintf("%s?uuid=%s&code=%s", Logic.Domain{}.GetAction(), req.Uuid, req.Code)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"action": domainInfo})
}

func (a otherAuth) Token(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
		Uuid string `json:"uuid"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	if req.Code == "" || req.Uuid == "" {
		Common.ApiResponse{}.Error(c, fmt.Sprintf("路由有误 code:%s uuid%s", req.Uuid, req.Uuid), gin.H{})
		return
	}

	var umap User.UserAuthMap
	Base.MysqlConn.Model(&umap).Where("cookie_uid = ?", req.Uuid).Find(&umap)

	service, err := Logic.Service{}.Get(req.Code)
	if err != nil {
		Common.ApiResponse{}.Error(c, req.Code+"客服不存在", gin.H{})
		return
	}

	// 如果用户不存在则创建新用户

	token := Common.Tools{}.EncodeToken(umap.UserId, "user", service.ServiceId, 0)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"token": token})
}

func (a otherAuth) ShowJoin(c *gin.Context) {
	c.HTML(http.StatusOK, "join.html", gin.H{})
}

func (a otherAuth) ShowAction(c *gin.Context) {
	c.HTML(http.StatusOK, "action.html", gin.H{})

}

func (a otherAuth) CodeToAction(c *gin.Context) {
	code := c.Query("code")
	var service Service.Service
	Base.MysqlConn.Model(&service).Where("code = ?", code).Find(&service)
	if service.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "客服不存在", gin.H{})
		return
	}

	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common.ApiResponse{}.Error(c, "客服已过期", gin.H{})
		return
	}

	if service.BindAction != "" {
		Common.ApiResponse{}.Success(c, "ok", gin.H{"action": service.BindAction})
		return
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"action": Logic.Domain{}.GetAction()})
}
