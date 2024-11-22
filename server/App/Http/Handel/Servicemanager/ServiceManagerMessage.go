package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Constant"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/Service"
	"server/App/Model/ServiceManager"
	"server/Base"
	"time"
)

type ServiceManagerMessage struct{}

// @summary 快捷回复管理-新增消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param title query string true "标题"
// @Param content query string true "标题内容"
// @Param type query string true "内容类型 text image"
// @Router /service/manager/message/add [post]
func (ServiceManagerMessage) Add(c *gin.Context) {
	var req Request.ServiceManagerMessageAddOrUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}
	req.AddServiceId = 0
	req.UpdateTime = time.Now()
	req.Status = "agree"
	req.ServiceManagerId = Common.Tools{}.GetRoleId(c)

	var model Response.ServiceManagerMessage
	model.Type = req.Type
	model.Content = req.Content
	model.Title = req.Title
	model.AddServiceId = 0
	model.UpdateTime = time.Now()
	model.CreateTime = time.Now()
	model.Status = "agree"
	model.ServiceManagerId = Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Create(&model)

	var serviceList []Service.Service
	Base.MysqlConn.Find(&serviceList, "service_manager_id = ?", req.ServiceManagerId)
	for _, item := range serviceList {
		Base.MysqlConn.Create(&Service.ServiceMessage{
			ServiceId:  item.ServiceId,
			MsgType:    "text",
			MsgInfo:    req.Content,
			Status:     "enable",
			Type:       "quick_reply",
			CreateTime: time.Now(),
		})
	}

	Common.ApiResponse{}.Success(c, "添加成功", gin.H{})
}

// @summary 快捷回复管理-获取消息列表
// @tags 客服后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/message/list [post]
func (ServiceManagerMessage) List(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	ServiceManagerId := Common.Tools{}.GetRoleId(c)

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Model(&ServiceManager.ServiceManagerMessage{}).Where("service_manager_id = ?", ServiceManagerId).Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.ServiceManagerMessage

	Base.MysqlConn.Where("service_manager_id = ?", ServiceManagerId).Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Find(&list)

	var listToTime []Response.ServiceManagerMessageTimeToString
	for _, v := range list {
		listToTime = append(listToTime, Response.ServiceManagerMessageTimeToString{
			Id:               v.Id,
			Title:            v.Title,
			Content:          v.Content,
			Type:             v.Type,
			ServiceManagerId: v.ServiceManagerId,
			AddServiceId:     v.AddServiceId,
			Status:           v.Status,
			CreateTime:       v.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:       v.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": listToTime}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 快捷回复管理-删除消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param id query string false "删除的指定ID"
// @Router /service/manager/message/delete [post]
func (ServiceManagerMessage) Delete(c *gin.Context) {
	var req Request.ServiceManagerMessageDelete
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}

	var info ServiceManager.ServiceManagerMessage
	Base.MysqlConn.Find(&info, "id = ?", req.Id)

	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Delete(&ServiceManager.ServiceManagerMessage{}, "id = ? and service_manager_id = ?", req.Id, ServiceManagerId)

	var serviceList []Service.Service
	Base.MysqlConn.Find(&serviceList, "service_manager_id = ?", ServiceManagerId)
	for _, item := range serviceList {
		Base.MysqlConn.Where("service_id = ? and msg_info = ?", item.ServiceId, info.Content).Delete(&Service.ServiceMessage{})
	}

	Common.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

// @summary 快捷回复管理-修改消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param id query string true "消息ID"
// @Param title query string true  "标题"
// @Param content query string true "标题内容"
// @Param type query string true "内容类型 text image"
// @Router /service/manager/message/update [post]
func (ServiceManagerMessage) Update(c *gin.Context) {
	var req Request.ServiceManagerMessageAddOrUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}

	var model ServiceManager.ServiceManagerMessage
	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&model).Where("id = ? and service_manager_id = ?", req.Id, ServiceManagerId).Find(&model)

	var serviceList []Service.Service
	Base.MysqlConn.Find(&serviceList, "service_manager_id = ?", ServiceManagerId)

	for _, item := range serviceList {
		Base.MysqlConn.Model(&Service.ServiceMessage{}).Where("service_id = ? and msg_info = ?", item.ServiceId, model.Content).Updates(
			gin.H{"msg_info": req.Content},
		)
	}

	model.Type = req.Type
	model.Content = req.Content
	model.Title = req.Title
	model.AddServiceId = 0
	model.UpdateTime = Constant.SystemTime{}
	model.CreateTime = Constant.SystemTime{}
	model.Status = "agree"
	model.ServiceManagerId = Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&model).Where("id = ? and service_manager_id = ?", req.Id, ServiceManagerId).Updates(model)

	Common.ApiResponse{}.Success(c, "修改成功", gin.H{})
}
