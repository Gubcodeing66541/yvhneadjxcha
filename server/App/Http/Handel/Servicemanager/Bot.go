package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/ServiceManager"
	"server/Base"
	"time"
)

type Bot struct{}

// @summary 机器人管理-获取机器人信息
// @tags 客服后台
// @Param token header string true "认证token"
// @Router /service/manager/bot/info [post]
func (Bot) Info(c *gin.Context) {
	var bot ServiceManager.ServiceManagerBot
	id := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Find(&bot, "service_manager_id = ?", id)
	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"info": bot})
}

// @summary 机器人管理-修改机器人信息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param status query string false "状态 run启动 stop停止"
// @Param head query string false "头像"
// @Param hello query string false "打招呼信息"
// @Router /service/manager/bot/update_info [post]
func (Bot) UpdateInfo(c *gin.Context) {
	var req Request.BotUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整信息", gin.H{})
		return
	}
	Id := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&ServiceManager.ServiceManagerBot{}).Where("service_manager_id = ?", Id).Updates(req)
	Common.ApiResponse{}.Success(c, "修改成功", gin.H{})
}

// @summary 机器人管理-获取机器人消息列表
// @tags 客服后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/bot/list [post]
func (Bot) List(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Model(&ServiceManager.ServiceManagerBotMessage{}).Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	var list []ServiceManager.ServiceManagerBotMessage

	Base.MysqlConn.Where("service_manager_id = ?", ServiceManagerId).Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Find(&list)
	var BotMessageResp []Response.ServiceManagerBotMessageResp
	for _, v := range list {
		BotMessageResp = append(BotMessageResp, Response.ServiceManagerBotMessageResp{
			Id:               v.Id,
			ServiceManagerId: v.ServiceManagerId,
			Problem:          v.Problem,
			Answer:           v.Answer,
			CreateTime:       v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": BotMessageResp}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 机器人管理-新增机器人消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param problem query string false "问题"
// @Param answer query string false "回答"
// @Router /service/manager/bot/add [post]
func (Bot) Add(c *gin.Context) {
	var req ServiceManager.ServiceManagerBotMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}
	req.ServiceManagerId = Common.Tools{}.GetRoleId(c)
	req.CreateTime = time.Now()
	Base.MysqlConn.Create(&req)
	Common.ApiResponse{}.Success(c, "添加成功", gin.H{})
}

// @summary 机器人管理-修改机器人消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param id query string false "问题id"
// @Param problem query string false "问题"
// @Param answer query string false "回答"
// @Router /service/manager/bot/update [post]
func (Bot) Update(c *gin.Context) {
	var req ServiceManager.ServiceManagerBotMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}
	req.ServiceManagerId = Common.Tools{}.GetRoleId(c)
	req.CreateTime = time.Now()
	Base.MysqlConn.Where("id = ? and service_manager_id = ?", req.Id, req.ServiceManagerId).Updates(&req)
	Common.ApiResponse{}.Success(c, "修改成功", gin.H{})
}

// @summary 机器人管理-删除机器人消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param id query string false "消息id"
// @Router /service/manager/bot/delete [post]
func (Bot) Delete(c *gin.Context) {
	var req ServiceManager.ServiceManagerBotMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请选择需要删除的数据", gin.H{})
		return
	}
	req.ServiceManagerId = Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Delete(&ServiceManager.ServiceManagerBotMessage{}, "id = ? and service_manager_id = ?", req.Id, req.ServiceManagerId)
	Common.ApiResponse{}.Success(c, "删除成功", gin.H{})
}
