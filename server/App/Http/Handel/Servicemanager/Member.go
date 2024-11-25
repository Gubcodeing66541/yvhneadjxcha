package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common2 "server/App/Model/Common"
	"server/App/Model/Service"
	ServiceManager2 "server/App/Model/ServiceManager"
	"server/Base"
	"strings"
	"time"
)

type Member struct{}

// @summary 账号管理-创建账号
// @tags 客服后台
// @Param token header string true "认证token"
// @Param service_name query string false "客服名称"
// @Param day query int false "开通天数"
// @Router /service/manager/member/create [post]
func (Member) Create(c *gin.Context) {
	var req Request.MemberCreateService
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	// 检测冻结
	serviceManager := Logic.ServiceManager{}.Get(Common.Tools{}.GetRoleId(c))
	if serviceManager.Status == "no_use" {
		Common.ApiResponse{}.Error(c, "您的账号已被冻结", gin.H{"req": req})
		return
	}
	// 检测余额是否足够
	pay := Logic.Order{}.GerAmount()
	if serviceManager.Account < pay*req.Day {
		Common.ApiResponse{}.Error(c, "余额不足", gin.H{"req": req})
		return
	}

	// 可用域名检测
	//num := Logic.Domain{}.GetNoUsePrivateNum()
	//if num == 0 {
	//	Common.ApiResponse{}.Error(c, "可用私有域名不足", gin.H{"req": req})
	//	return
	//}
	// 注册
	serviceId, member := Logic.Service{}.CreateByServiceManager(Common.Tools{}.GetRoleId(c), req.ServiceName, req.Day)

	// 生成快捷回复
	var message []ServiceManager2.ServiceManagerMessage
	Base.MysqlConn.Find(&message, "service_manager_id = ?", serviceManager.ServiceManagerId)
	for _, item := range message {
		Base.MysqlConn.Create(&Service.ServiceMessage{
			ServiceId: serviceId, MsgType: "text", MsgInfo: item.Content, Status: "enable", Type: "quick_reply",
			CreateTime: time.Now(),
		})
	}

	// 续费
	_ = Logic.Service{}.RenewalByServiceManager(Common.Tools{}.GetRoleId(c), member, req.Day, "create_service")
	Common.ApiResponse{}.Success(c, "账号开通成功", gin.H{"member": member, "req": req})
}

// @summary 账号管理-批量创建账号
// @tags 客服后台
// @Param token header string true "认证token"
// @Param service_number query int false "开通数量"
// @Param day query int false "开通天数"
// @Router /service/manager/member/create_list [post]
func (Member) CreateList(c *gin.Context) {
	var req Request.MemberCreateServiceList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	serviceManager := Logic.ServiceManager{}.Get(Common.Tools{}.GetRoleId(c))

	// 冻结检测
	if serviceManager.Status == "no_use" {
		Common.ApiResponse{}.Error(c, "您的账号已被冻结", gin.H{"req": req})
		return
	}

	// 检测余额是否足够
	pay := Logic.Order{}.GerAmount()
	if serviceManager.Account < pay*req.ServiceNumber*req.Day {
		Common.ApiResponse{}.Error(c, "余额不足", gin.H{"req": req})
		return
	}

	// 域名数量检测
	//num := Logic.Domain{}.GetNoUsePrivateNum()
	//if num < req.ServiceNumber {
	//	Common.ApiResponse{}.Error(c, "可用私有域名不足", gin.H{"req": req})
	//	return
	//}

	//快捷回复
	var message []ServiceManager2.ServiceManagerMessage
	Base.MysqlConn.Find(&message, "service_manager_id = ?", serviceManager.ServiceManagerId)

	var memberList []string
	for i := 0; i < req.ServiceNumber; i++ {

		// 创建账号
		serviceId, member := Logic.Service{}.CreateByServiceManager(Common.Tools{}.GetRoleId(c), "小客服", req.Day)

		// 生成快捷回复
		for _, item := range message {
			Base.MysqlConn.Create(Service.ServiceMessage{
				ServiceId: serviceId, MsgType: "text", MsgInfo: item.Content, Status: "enable", Type: "quick_reply",
				CreateTime: time.Now(),
			})
		}

		// 续费
		_ = Logic.Service{}.RenewalByServiceManager(Common.Tools{}.GetRoleId(c), member, req.Day, "create_service")
		memberList = append(memberList, member)
	}
	Common.ApiResponse{}.Success(c, "多账号开通成功", gin.H{"member_list": memberList})
}

// @summary 账号管理-获取账号列表
// @tags 客服后台
// @Param token header string true "认证token"
// @Param search query string false "搜索名字"
// @Param start_time query time false "开始时间"
// @Param end_time query time false "结束时间"
// @Param page query int false "指定页"
// @Param offset query int true "分页数量"
// @Router /service/manager/member/list [post]
func (Member) List(c *gin.Context) {
	var pageReq Request.PageLimitByDate
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
	}

	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	where := Base.MysqlConn.Model(&Service.Service{}).Where("service_manager_id = ?", ServiceManagerId)

	// 条件对比1
	if pageReq.Search != "" {
		where = where.Where("name like ? ", "%"+pageReq.Search+"%")
	}

	// 条件对比1
	if pageReq.Username != "" {
		where = where.Where("username like ?", "%"+pageReq.Username+"%")
	}

	// 条件对比2
	if pageReq.StartTime != "" {
		where = where.Where("create_time BETWEEN ? AND ?", pageReq.StartTime, pageReq.EndTime)
	}

	// 计算分页和总数
	var allCount int
	where.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.ServiceList

	where.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)

	for key, _ := range list {
		serviceSocketKey := Common.Tools{}.GetServiceWebSocketId(list[key].ServiceId)
		list[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(serviceSocketKey)
		Base.MysqlConn.Model(&Service.ServiceRoom{}).Where("service_id = ?", list[key].ServiceId).Count(&list[key].UserCnt)
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
	Common.ApiResponse{}.Success(c, "ok", res)
}

// @summary 账号管理-修改账号信息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param service_id query int false "客服ID"
// @Param name query string false "客服名称"
// @Router /service/manager/member/update [post]
func (Member) Update(c *gin.Context) {
	var req Request.MemberCreateServiceUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}
	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Service.Service{}).Where("service_id = ? and service_manager_id = ?", req.ServiceId, ServiceManagerId).Updates(req)
	Common.ApiResponse{}.Success(c, "修改成功", gin.H{})
}

// @summary 账号管理-删除账号
// @tags 客服后台
// @Param token header string true "认证token"
// @Param service_id query int false "删除的客服ID"
// @Router /service/manager/member/delete [post]
func (Member) Delete(c *gin.Context) {
	var req Request.MemberCreateServiceId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}

	ServiceManagerId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Common2.Domain{}).Where("bind_service_id = ?", req.ServiceId).Updates(gin.H{"bind_service_id": 0})
	Base.MysqlConn.Delete(&Service.Service{}, "service_id = ? and service_manager_id = ?", req.ServiceId, ServiceManagerId)
	Base.MysqlConn.Delete(&Service.ServiceAuth{}, "service_id = ? and service_manager_id = ?", req.ServiceId, ServiceManagerId)
	Common.ApiResponse{}.SendMsgToService(req.ServiceId, "out_login", gin.H{"message": "您的账号已被删除"})
	Common.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

// @summary 账号管理-续费
// @tags 客服后台
// @Param token header string true "认证token"
// @Param username query string false "客服ID"
// @Param day query int false "续费的天数"
// @Router /service/manager/member/renewal [post]
func (Member) Renewal(c *gin.Context) {
	var req Request.MemberServiceRenewal
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}

	serviceManagerId := Common.Tools{}.GetRoleId(c)
	err = Logic.Service{}.RenewalByServiceManager(serviceManagerId, req.Username, int(req.Day), "renew_service")
	if err != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}
	Common.ApiResponse{}.Success(c, "续费成功", gin.H{"req": req})
}

// @summary 账号管理-批量续费
// @tags 客服后台
// @Param token header string true "认证token"
// @Param username_list query string false "客服账号列表 \r 分割"
// @Param day query int false "续费的天数"
// @Router /service/manager/member/renewal_all [post]
func (Member) RenewalAll(c *gin.Context) {
	var req Request.MemberServiceRenewalList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请填写完整", gin.H{})
		return
	}

	// 将字符串 分割成 字符串数组
	// 参数：要拼接的字符串，分割的内容
	memberList := strings.Split(req.UsernameList, "\n")
	for _, username := range memberList {
		if username == "" {
			continue
		}
		serviceManagerId := Common.Tools{}.GetRoleId(c)
		err = Logic.Service{}.RenewalByServiceManager(serviceManagerId, username, int(req.Day), "renew_service")
		if err != nil {
			Common.ApiResponse{}.Error(c, err.Error(), gin.H{"username": username})
			return
		}
	}
	Common.ApiResponse{}.Success(c, "续费成功", gin.H{"req": req})
}
