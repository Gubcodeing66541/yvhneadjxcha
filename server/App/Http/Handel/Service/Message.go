package Service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Tel"
	Message2 "server/App/Model/Message"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Message struct{}

// @summary 聊天-发送消息给用户
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "发送人ID"
// @Param type query string true "消息类型"
// @Param content query string true "消息内容"
// @Router /service/message/send_to_user [post]
func (Message) SendToUser(c *gin.Context) {
	var req Request.ServiceSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 || req.Content == "" {
		Common.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	RoleId := Common.Tools{}.GetRoleId(c)
	err = Logic.Message{}.SendToUser(RoleId, req.UserId, req.Type, req.Content, true)
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{"erq": req})
	}

	// 返回OK信息
	Common.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})

}

// @summary 聊天-群发消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query array false "查询的用户ID数组"
// @Param type query string true "消息类型"
// @Param content query string true "消息内容"
// @Router /service/message/send_all [post]
func (Message) SendAll(c *gin.Context) {
	var req Request.ServiceSendMessageGroup
	userList := c.PostFormArray("user_id")

	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" || req.Type == "" {
		Common.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req, "userList": userList})
		return
	}

	// 开启协程循环发送
	RoleId := Common.Tools{}.GetRoleId(c)
	go func(RoleId int, req Request.ServiceSendMessageGroup) {
		// 查询非黑名单用户列表
		var i = 0
		for _, Item := range req.UserId {
			i++
			err = Logic.Message{}.SendToUser(RoleId, Item, req.Type, req.Content, true)
			if err != nil {
				fmt.Println("send Error", err.Error())
			}
			if i >= 50 {
				i = 0
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(RoleId, req)

	Common.ApiResponse{}.Success(c, "群发成功", gin.H{"erq": req})

	return
}

// @summary 聊天-获取聊天记录
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query string false "查询的用户ID"
// @Param page query int true "指定页"
// @Param offset query int true "指定每页数量"
// @Router /service/message/list [post]
func (Message) List(c *gin.Context) {
	var req Request.MsgList
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	var message []Tel.Message
	tel := Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ? and is_del = 0", roleId, req.UserId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel = Base.MysqlConn.Raw("select * from (select * from messages where service_id = ? and user_id = ?  and is_del = 0 order by id desc limit ? offset ? )  t order by id asc",
		roleId, req.UserId, req.Offset, (req.Page-1)*req.Offset)
	tel.Scan(&message)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": message})
}

func (Message) Update(c *gin.Context) {
	var req Request.UpdateServiceDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", roleId).Updates(req)
	Logic.Service{}.ClearCache(roleId)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"REQ": req})
}

// @summary 聊天-消息撤回
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query string false "查询的用户ID"
// @Param id query int true "指定的消息ID"
// @Router /service/message/remove_msg [post]
func (Message) RemoveMessage(c *gin.Context) {
	var req Request.RemoveMsg
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "填写完整信息", gin.H{})
		return
	}

	//逻辑删除消息
	//Base.MysqlConn.Unscoped().Delete(&Message2.Message{}, "id = ?", req.Id)
	Base.MysqlConn.Model(&Message2.Message{}).Where("id = ?", req.Id).Updates(
		&Message2.Message{IsDel: 1})

	//是不是最近消息

	var lateMessage Service2.ServiceRoom
	Base.MysqlConn.Find(&lateMessage, "late_id = ? ", req.Id)
	RoleId := Common.Tools{}.GetServiceId(c)
	if req.Id == lateMessage.LateId {
		Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(gin.H{"late_msg": "你撤回了一条消息", "late_type": "text"})
	}

	Common.ApiResponse{}.SendMsgToService(RoleId, "remove", req)
	Common.ApiResponse{}.SendMsgToUser(req.UserId, "remove", req)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}

// @summary 聊天-消息清空
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query string false "查询的用户ID"
// @Router /service/message/clear_message [post]
func (Message) ClearMessage(c *gin.Context) {

	var req Request.RemoveMsg
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "填写完整信息", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetServiceId(c)
	//Base.MysqlConn.Delete(&Message2.Message{}, "service_id = ? and user_id = ?", RoleId, req.UserId)

	Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(
		&Message2.Message{IsDel: 1})

	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", RoleId, req.UserId).Updates(gin.H{
		"LateId":         0,
		"LateType":       "",
		"LateMsg":        "",
		"LateRole":       "",
		"LateUserReadId": 0,
		"UserNoRead":     0,
		"ServiceNoRead":  0,
	})
	Common.ApiResponse{}.SendMsgToService(RoleId, "clear", req)
	Common.ApiResponse{}.SendMsgToUser(req.UserId, "clear", req)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}
