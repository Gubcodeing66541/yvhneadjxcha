package Manage

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Service struct{}

// 创建用户号
func (Service) ServiceCreate(c *gin.Context) {
	var req Request.CreateServiceDay
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "创建失败", gin.H{})
		return
	}

	//todo 创建类型 auth还是push
	serviceId, member, password := Logic.Service{}.Create()
	res := gin.H{"member": member, "password": password, "service_id": serviceId}

	// 创建的时候执行续费逻辑
	_ = Logic.Service{}.Renewal(serviceId, req.Day)
	Common.ApiResponse{}.Success(c, "注册成功", res)
}

// 批量创建用户号
func (Service) ServiceBachCreate(c *gin.Context) {
	var req Request.CreateServiceDay
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数不全", gin.H{})
		return
	}

	//批量创建并且续费
	service := Logic.Service{}.BachCreate(req)
	for i := 0; i < len(service); i++ {
		err2 := Logic.Service{}.Renewal(service[i].ServiceId, req.Day)
		if err2 != nil {
			Common.ApiResponse{}.Success(c, "批量续费失败", gin.H{})
		}
	}
	res := gin.H{"service": service}
	Common.ApiResponse{}.Success(c, "注册成功", res)
}

func (Service) ServiceList(c *gin.Context) {
	var req Request.GerServiceList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	services := Logic.Service{}.List(req)
	Common.ApiResponse{}.Success(c, "获取客服列表", gin.H{"services": services})

}

// 获取客服域名
func (Service) GetServiceDomain(c *gin.Context) {
	var req Request.GetServiceDomain
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	domain := Logic.Service{}.GetServiceDomain(req.ServiceId)
	Common.ApiResponse{}.Success(c, "获取客服域名成功", gin.H{"domain": domain})
}

// 续费
func (Service) Renewal(c *gin.Context) {
	var req Request.Renewal
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.Service{}.Renewal(req.ServiceId, req.Day)
	if err2 != nil {
		Common.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}

	Common.ApiResponse{}.Success(c, "续费成功", gin.H{})
}

// 换绑
func (Service) ChangeBindDomain(c *gin.Context) {
	var req Request.ChangeBindDomain
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	domain := Logic.Domain{}.GetServiceBind(req.ServiceId)

	//解绑 换绑
	Logic.Domain{}.EnableDisable(domain.Id, "un_enable")
}

// 查询客服的订单
func (Service) GetServiceOrder(c *gin.Context) {
	var req Request.GetServiceOrder
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	orders := Logic.Order{}.GetByServiceId(req.ServiceId)
	Common.ApiResponse{}.Success(c, "获取客服订单", gin.H{"orders": orders})
}

// 查询客服的订单详情
func (Service) GetServiceOrderInfo(c *gin.Context) {
	var req Request.GetServiceOrderInfo
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	orders := Logic.Order{}.QueryById(req.OrderId)
	Common.ApiResponse{}.Success(c, "获取订单详情", gin.H{"orders": orders})
}

// 绑定域名
func (Service) BindDomain(c *gin.Context) {
	var req Request.BindDomain
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.Domain{}.Bind(req.ServiceId)
	if err2 != nil {
		Common.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common.ApiResponse{}.Success(c, "绑定域名", nil)
}

// @summary 在线统计
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Router /manager/count [post]
func (Service) Count(c *gin.Context) {
	sql := "select " +
		"(select sum(renew) as renew_service_manager from service_manager_renew_recorders where reason = 'renew_service_manager') as renew_service_manager," +
		"(select sum(account) as all_account from service_managers) as all_account," +
		"(select sum(renew) as all_pay from service_manager_renew_recorders where reason != 'renew_service_manager') as all_pay," +
		"(select sum(renew) as all_pay from service_manager_renew_recorders where reason != 'renew_service_manager' and create_time >= ?) as today_pay"
	var count Response.ManagerPayCont
	Base.MysqlConn.Raw(sql, time.Now().Format("2006-01-02")).Scan(&count)

	serviceOnlineNum := 0
	onlineAllNum := Base.WebsocketHub.GetOnlineCount() - serviceOnlineNum

	serviceIds, user := Common.Socket{}.GetAll()
	serviceCount, userCount := len(serviceIds), len(user)

	var services []Service2.Service
	Base.MysqlConn.Where("service_id in (?)", serviceIds).Find(&services)

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{
		"service_on_line_num": serviceOnlineNum,
		"online_all_num":      onlineAllNum,
		"count":               count,
		"socket": gin.H{
			"service_count":  serviceCount,
			"user_count":     userCount,
			"socket_count":   Base.WebsocketHub.GetOnlineCount(),
			"online_service": services,
		},
	})
}

// @summary 在线统计-客服列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/count_service_list [post]
func (Service) CountServiceList(c *gin.Context) {
	var pageReq Request.PageLimitByDate
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	serviceIds, _ := Common.Socket{}.GetAll()

	var services []Response.ServiceListBySocket
	tel := Base.MysqlConn.Raw("select services.*,service_managers.member from services "+
		"left join service_managers on services.service_manager_id = service_managers.service_manager_id "+
		" where service_id in (?)", serviceIds)

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Raw("select count(*) from services where service_id in (?)", serviceIds).Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&services)

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{
		"count": allCount, "page": allPage, "current_page": pageReq.Page,
		"services": services,
	})

}
