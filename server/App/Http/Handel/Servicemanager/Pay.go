package Servicemanager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/ServiceManager"
	"server/Base"
)

type Pay struct{}

func (Pay) Create(c *gin.Context) {

}

// @summary 账单记录
// @tags 客服后台
// @Param token header string true "认证token"
// @Param search query string false "查询的账户"
// @Param reason query string false "类型"
// @Param start_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Param page query int true "指定页"
// @Param offset query int true "指定每页数量"
// @Router /service/manager/pay/recorder [post]
func (Pay) Recorder(c *gin.Context) {
	var req Request.PayRecorder
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "用户不存在", gin.H{})
		return
	}

	where := fmt.Sprintf("service_manager_id = %d", Common.Tools{}.GetRoleId(c))

	// 查询总充值和总新增账号金额以及总计续费金额
	var payCount Response.PayCount
	sql := "select (select sum(renew) as money from service_manager_renew_recorders where service_manager_id = ? and reason = 'renew_service_manager') as renew_service_manager," +
		"(select sum(renew) as money from service_manager_renew_recorders where service_manager_id = ? and reason = 'renew_service')  as renew_service," +
		"(select sum(renew) as money from service_manager_renew_recorders where service_manager_id = ? and reason = 'create_service') as create_service"
	RoleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Raw(sql, RoleId, RoleId, RoleId).Scan(&payCount)

	if req.Search != "" {
		where += fmt.Sprintf(" and member = '%s'", req.Search)
	}

	if req.ServiceManagerId != 0 {
		where += fmt.Sprintf(" and service_manager_id = %d", req.ServiceManagerId)
	}

	if req.Reason != "" {
		where += fmt.Sprintf(" and reason = '%s'", req.Reason)
	}
	if req.StartTime != "" {
		where += fmt.Sprintf(" and create_time >= '%s'", req.StartTime)
	}

	if req.EndTime != "" {
		where += fmt.Sprintf(" and create_time <= '%s'", req.EndTime)
	}

	var message []ServiceManager.ServiceManagerRenewRecorder
	tel := Base.MysqlConn.Model(&ServiceManager.ServiceManagerRenewRecorder{}).Where(where)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel.Offset((req.Page - 1) * req.Offset).Limit(req.Offset).Find(&message)

	var recorder []Response.ServiceManagerRenewRecorder
	for _, messagedata := range message {
		renewRecorder := Response.ServiceManagerRenewRecorder{
			Id:                   messagedata.Id,
			OrderId:              messagedata.OrderId,
			ServiceManagerId:     messagedata.ServiceManagerId,
			ServiceManagerMember: messagedata.ServiceManagerMember,
			Member:               messagedata.Member,
			ServiceId:            messagedata.ServiceId,
			OldAccount:           messagedata.OldAccount,
			Account:              messagedata.Account,
			Renew:                messagedata.Renew,
			Reason:               messagedata.Reason,
			PayType:              messagedata.PayType,
			CreateTime:           messagedata.CreateTime.Format("2006-01-02 15:04:05"),
		}
		recorder = append(recorder, renewRecorder)
	}

	Common.ApiResponse{}.Success(c, "OK", gin.H{"count": allCount, "page": allPage, "pay_count": payCount, "current_page": req.Page, "list": recorder})
}
