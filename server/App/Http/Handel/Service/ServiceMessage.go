package Service

import (
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	Service2 "server/App/Model/Service"
	"server/Base"
)

type ServiceMessage struct{}

// @summary 消息管理-创建消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param msg_info query string true "消息内容"
// @Param msg_type query string true "消息类型"
// @Param type query string true "创建类型 hello打招呼 group群发消息 quick_reply快捷回复"
// @Router /service/service_message/create [post]
func (ServiceMessage) Create(c *gin.Context) {
	var req Request.CreateServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Logic.ServiceMessage{}.Create(serviceId, req)
	Common2.ApiResponse{}.Success(c, "创建成功", gin.H{})
}

// @summary 消息管理-删除消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param id query int true "删除的ID"
// @Router /service/service_message/delete [post]
func (ServiceMessage) Delete(c *gin.Context) {
	var req Request.DeleteServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Logic.ServiceMessage{}.Delete(req.Id, req.ServiceId)
	Common2.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

// @summary 消息管理-修改消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param id query string true "消息ID"
// @Param msg_info query string true "消息内容"
// @Param msg_type query string true "消息类型"
// @Param status query string true "启用状态"
// @Param type query string true "创建类型 hello打招呼 group群发消息 quick_reply快捷回复"
// @Router /service/service_message/update [post]
func (ServiceMessage) Update(c *gin.Context) {
	var req Request.UpdateServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	Logic.ServiceMessage{}.Update(req.Id, req.MsgType, req.MsgInfo, req.ServiceId, req.Status)
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 消息管理-消息列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param type query string true "类型"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/service_message/list [post]
func (ServiceMessage) List(c *gin.Context) {
	var pageReq Request.ListServiceMessage
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)

	tel := Base.MysqlConn.Model(&Service2.ServiceMessage{}).Where("service_id = ? and type = ?", ServiceId, pageReq.Type)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Service2.ServiceMessage
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 消息管理-获取单条消息详细
// @tags 客服系统
// @Param token header string true "认证token"
// @Param id query string true "消息ID"
// @Router /service/service_message/get [post]
func (ServiceMessage) GetById(c *gin.Context) {
	var req Request.GetByIdServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	req.ServiceId = Common2.Tools{}.GetServiceId(c)
	serviceMessage := Logic.ServiceMessage{}.GetById(req.ServiceId, req.Id)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{"serviceMessage": serviceMessage})
}

// @summary 消息管理-位置交换
// @tags 客服系统
// @Param token header string true "认证token"
// @Param from query int true "消息ID1"
// @Param to query int true "消息ID2"
// @Router /service/service_message/swap [post]
func (ServiceMessage) Swap(c *gin.Context) {
	var req Request.SwapServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)
	from := Logic.ServiceMessage{}.GetById(ServiceId, req.From)
	to := Logic.ServiceMessage{}.GetById(ServiceId, req.To)
	from.Id, to.Id = to.Id, from.Id
	Base.MysqlConn.Save(&from)
	Base.MysqlConn.Save(&to)
	Common2.ApiResponse{}.Success(c, "修改成功", gin.H{"from": from, "to": to})
}
