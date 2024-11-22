package Servicemanager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Message2 "server/App/Model/Message"
	"server/Base"
)

type User struct{}

// @summary 用户列表
// @tags 客服后台
// @Param token header string true "认证token"
// @Param service_name query string false "客服名称"
// @Param user_name query string false "用户名称"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/users/list [post]
func (User) Users(c *gin.Context) {
	var req Request.ServiceRoomList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetRoleId(c)
	res := Logic.ServiceRoom{}.ListByServiceManager(RoleId, req)
	if req.IsClearStatus == 1 {
		Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), 0)
	}

	//序列化
	Common.ApiResponse{}.Success(c, "ok", res)
}

// @summary 聊天记录
// @tags 客服后台
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Param service_id query int true "客服ID"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/users/message [post]
func (User) Message(c *gin.Context) {
	var req Request.MsgList
	err := c.ShouldBind(&req)
	if err != nil || req.UserId == 0 {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	roleId := req.ServiceId
	var message []Message2.Message
	tel := Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", roleId, req.UserId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel = Base.MysqlConn.Raw("select * from (select * from messages where service_id = ? and user_id = ? order by id desc limit ? offset ? )  t order by id asc",
		roleId, req.UserId, req.Offset, (req.Page-1)*req.Offset)
	tel.Scan(&message)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": message})
}

// @summary 拉黑列表
// @tags 客服后台
// @Param token header string true "认证token"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param user_name query string false "用户名字"
// @Param service_name query string false "客服名字"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/users/black [post]
func (User) Black(c *gin.Context) {
	var req Request.BlackPage
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	ServiceNameWhere := ""
	if req.ServiceName != "" {
		ServiceNameWhere += fmt.Sprintf(" and name = '%s'", req.ServiceName)
	}

	UserNameWhere := ""
	if req.UserName != "" {
		UserNameWhere += fmt.Sprintf(" and users.user_name = '%s'", req.UserName)
	}

	if req.StartTime != "" {
		UserNameWhere += fmt.Sprintf(" and service_blacks.create_time >= '%s'", req.StartTime)
	}

	if req.EndTime != "" {
		UserNameWhere += fmt.Sprintf(" and service_blacks.create_time <= '%s'", req.EndTime)
	}

	sql := "select count(*) as cnt from service_blacks " +
		"left join users on service_blacks.user_id = users.user_id " +
		"left join services on service_blacks.service_id = services.service_id " +
		" where service_blacks.service_manager_id = ? " + ServiceNameWhere + " " + UserNameWhere
	tel := Base.MysqlConn.Raw(sql, Common.Tools{}.GetRoleId(c))

	var allCount Response.Counted
	tel.Find(&Response.Counted{}).Scan(&allCount)

	allPage := math.Ceil(float64(allCount.Cnt) / float64(req.Offset))

	var list []Response.ServiceBlack
	sql = "select * from service_blacks " +
		"left join users on service_blacks.user_id = users.user_id " +
		"left join services on service_blacks.service_id = services.service_id " +
		" where service_blacks.service_manager_id = ? " + ServiceNameWhere + UserNameWhere + " limit ? offset ? "
	tel = Base.MysqlConn.Raw(sql, Common.Tools{}.GetRoleId(c), req.Offset, (req.Page-1)*req.Offset).Scan(&list)
	Common.ApiResponse{}.Success(c, "OK", gin.H{"count": allCount.Cnt, "page": allPage, "current_page": req.Page, "list": list})
}
