package Manage

import (
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	ServiceManager2 "server/App/Model/ServiceManager"
	"server/Base"
	"time"
)

type ServiceManager struct{}

// @summary 账号-创建账号
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param count query int true "账号数量"
// @Router /manager/service_manager/create [post]
func (ServiceManager) Create(c *gin.Context) {
	var req Request.Count
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var list []gin.H
	for i := 0; i < req.Count; i++ {
		username, password := Logic.ServiceManager{}.Create(c, 0)
		list = append(list, gin.H{"username": username, "password": password})
	}

	Common2.ApiResponse{}.Success(c, "账号创建成功", gin.H{"list": list})
}

// @summary 账号-获取账号列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param search query string false "查询条件"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/service_manager/list [post]
func (ServiceManager) List(c *gin.Context) {
	var pageReq Request.PageLimitByDate
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	tel := Base.MysqlConn.Model(&ServiceManager2.ServiceManager{}).Select("service_managers.*,t.service_cnt")
	if pageReq.Search != "" {
		tel = tel.Where("member like ?", "%"+pageReq.Search+"%")
	}

	// 条件对比2
	if pageReq.StartTime != "" {
		tel = tel.Where("create_time BETWEEN ? AND ?", pageReq.StartTime, pageReq.EndTime)
	}

	t := "select service_manager_id,count(*) as service_cnt from services group by service_manager_id"
	tel = tel.Joins("left join (" + t + ") t on service_managers.service_manager_id = t.service_manager_id")

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.ServiceManagerDetail
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 账号-续费
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param service_manager_id query string true "指定ID"
// @Param account query int true "续费金额"
// @Router /manager/service_manager/renew [post]
func (ServiceManager) ReNew(c *gin.Context) {
	var req Request.Account
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请填写完整信息", gin.H{})
		return
	}
	Logic.ServiceManager{}.Renew(req.ServiceManagerId, req.Account, "renew_service_manager", "system", "")
	Common2.ApiResponse{}.Success(c, "充值成功", gin.H{"req": req})
}

// @summary 账号-获取客服关系后台的客服账号
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param service_manager_id query int false "需要查询的母账号ID"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/service_manager/get_service_list [post]
func (ServiceManager) GetServiceList(c *gin.Context) {
	var pageReq Request.ServiceManagerId
	err := c.ShouldBind(&pageReq)

	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整", gin.H{"err": err.Error()})
		return
	}

	tel := Base.MysqlConn.Model(&Service2.Service{}).Where("service_manager_id = ?", pageReq.ServiceManagerId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.ServiceList
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Order("id desc").Limit(pageReq.Offset).Scan(&list)
	for key, _ := range list {
		serviceSocketKey := Common2.Tools{}.GetServiceWebSocketId(list[key].ServiceId)
		list[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(serviceSocketKey)
		Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ?", list[key].ServiceId).Count(&list[key].UserCnt)
		status := "success"

		if list[key].IsActivate == 0 {
			status = "no_active"
		} else if time.Now().After(list[key].TimeOut.ToTime()) {
			status = "time_out"
		} else if list[key].Status == "no_use" {
			status = "no_use"
		}
		list[key].Status = status
	}

	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 账号-删除所有账号
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param service_manager_id query int false "需要删除的母账号ID"
// @Router /manager/service_manager/delete [post]
func (ServiceManager) Delete(c *gin.Context) {
	var pageReq Request.ServiceManagerId
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整", gin.H{"err": err.Error()})
		return
	}

	// 检测余额
	serviceManager := Logic.ServiceManager{}.Get(pageReq.ServiceManagerId)
	if serviceManager.Account != 0 {
		Common2.ApiResponse{}.Error(c, "账户还有余额，无法直接删除", gin.H{})
		return
	}

	// 解绑域名并推送下线消息
	var services []Service2.Service
	Base.MysqlConn.Find(&services, "service_manager_id = ?", pageReq.ServiceManagerId)

	// 删除账号
	Base.MysqlConn.Delete(&ServiceManager2.ServiceManager{}, "service_manager_id = ?", pageReq.ServiceManagerId)
	Base.MysqlConn.Delete(&Service2.Service{}, "service_manager_id = ?", pageReq.ServiceManagerId)
	Base.MysqlConn.Delete(&Service2.ServiceAuth{}, "service_manager_id = ?", pageReq.ServiceManagerId)

	for _, item := range services {
		Base.MysqlConn.Model(&Common.Domain{}).Where("bind_service_id = ?", item.ServiceId).Updates(gin.H{"bind_service_id": 0})
		Common2.ApiResponse{}.SendMsgToService(item.ServiceId, "out_login", gin.H{"message": "您的账号已被删除"})
		Base.MysqlConn.Delete(&Service2.ServiceRoom{}, "service_id = ?", item.ServiceId)
		Base.MysqlConn.Delete(&Service2.ServiceBlack{}, "service_id = ?", item.ServiceId)

	}

	Common2.ApiResponse{}.Success(c, "删除成功", gin.H{"req": pageReq})
}

// @summary 账号-修改指定账号密码
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param service_manager_id query int false "账号ID"
// @Param name query int false "修改名字"
// @Param password query int false "修改密码"
// @Router /manager/service_manager/reset_password [post]
func (ServiceManager) ResetPassword(c *gin.Context) {
	var pageReq Request.ServiceManagerReset
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整", gin.H{"err": err.Error()})
		return
	}

	if pageReq.Password != "" {
		Base.MysqlConn.Model(&ServiceManager2.ServiceManager{}).
			Where("service_manager_id = ?", pageReq.ServiceManagerId).
			Updates(gin.H{"password": pageReq.Password, "name": pageReq.Name})
	} else {
		Base.MysqlConn.Model(&ServiceManager2.ServiceManager{}).
			Where("service_manager_id = ?", pageReq.ServiceManagerId).
			Updates(gin.H{"name": pageReq.Name})
	}

	Common2.ApiResponse{}.Success(c, "修改成功", gin.H{"req": pageReq})
}

// @summary 账号-冻结指定账号
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param service_manager_id query int false "账号ID"
// @Router /manager/service_manager/ban [post]
func (s ServiceManager) Ban(c *gin.Context) {
	var pageReq Request.ServiceManagerId
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整", gin.H{"err": err.Error()})
		return
	}

	if pageReq.Status == "success" {
		// ban所有的账号
		Base.MysqlConn.Model(ServiceManager2.ServiceManager{}).Where("service_manager_id = ?", pageReq.ServiceManagerId).
			Updates(gin.H{"status": "success"})

		Base.MysqlConn.Model(Service2.Service{}).Where("service_manager_id = ?", pageReq.ServiceManagerId).
			Updates(gin.H{"status": "success"})

		Common2.ApiResponse{}.Success(c, "解除冻结成功", gin.H{"req": pageReq})
		return
	}

	// ban所有的账号
	Base.MysqlConn.Model(ServiceManager2.ServiceManager{}).Where("service_manager_id = ?", pageReq.ServiceManagerId).
		Updates(gin.H{"status": "no_use"})

	Base.MysqlConn.Model(Service2.Service{}).Where("service_manager_id = ?", pageReq.ServiceManagerId).
		Updates(gin.H{"status": "no_use"})

	var serviceList []Service2.Service
	Base.MysqlConn.Find(&serviceList, "service_manager_id = ?", pageReq.ServiceManagerId)
	for _, item := range serviceList {
		Common2.ApiResponse{}.SendMsgToService(item.ServiceId, "out_login", gin.H{"message": "您的账号已被冻结"})
	}

	Common2.ApiResponse{}.Success(c, "冻结成功", gin.H{"req": pageReq})
}
