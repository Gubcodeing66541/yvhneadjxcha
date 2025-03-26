package Manage

import (
	"fmt"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/Common"
	"server/Base"

	"github.com/gin-gonic/gin"
)

type Domain struct {
}

// @summary 域名-列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param type query string true "类型 action 罗迪 "
// @Param domain query string false "搜索域名"
// @Param username query string fasle "搜索账号"
// @Param is_bind_service query int fasle "绑定状态 0所有 1是绑定  2未绑定"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/domain/list [post]
func (Domain) List(c *gin.Context) {
	var pageReq Request.DomainListLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	tel := Base.MysqlConn.Table("domains").
		Select("domains.*, services.*,service_managers.service_manager_id,service_managers.name as service_manager_name,service_managers.member as service_manager_member").
		Joins("left JOIN services ON services.service_id = domains.bind_service_id").
		Joins("left JOIN service_managers ON services.service_manager_id = service_managers.service_manager_id")

	tel = tel.Where("domains.type = ? ", pageReq.Type)

	if pageReq.Domain != "" {
		tel = tel.Where("domains.domain like ?", "%"+pageReq.Domain+"%")
	}

	if pageReq.Username != "" {
		tel = tel.Where("services.username like ? ", "%"+pageReq.Username+"%")
	}

	if pageReq.IsBindService == 1 {
		tel = tel.Where("bind_service_id != 0")
	}

	if pageReq.IsBindService == 2 {
		tel = tel.Where("bind_service_id = 0")
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.RespDomainList
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list, "page_req": pageReq, "limit": pageReq.Offset, "offset": pageReq.Page - 1*pageReq.Offset}
	Common2.ApiResponse{}.Success(c, "获取成功", res)

}

// @summary 域名-通过ID查询域名信息
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param domain_id query string true "域名ID"
// @Router /manager/domain/query_by_id [post]
func (Domain) QueryById(c *gin.Context) {
	var req Request.QueryById
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	domain := Logic.Domain{}.QueryById(req.DomainId)
	if domain.Id == 0 {
		Common2.ApiResponse{}.Error(c, "未查询到该域名", gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "id查询域名", gin.H{"domains": domain})
}

// @summary 域名-删除域名
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param domain_id query string true "域名ID"
// @Router /manager/domain/delete [post]
func (Domain) Delete(c *gin.Context) {
	var req Request.DomainDelete
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.Domain{}.Delete(req.DomainId)
	if err2 != nil {
		Common2.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "删除域名", gin.H{})
}

// @summary 域名-修改域名
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query string true "域名ID"
// @Param domain query string true "域名"
// @Param type query string true "域名类型 private action "
// @Param status query string true "状态 enable 启动 un_enable 禁用"
// @Router /manager/domain/update [post]
func (Domain) Update(c *gin.Context) {
	var req Request.DomainUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.Domain{}.Update(req.Id, req.Domain, req.Type, req.Status)
	if err2 != nil {
		Common2.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common2.ApiResponse{}.Success(c, "修改域名", gin.H{})
}

// @summary 域名-添加域名
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param domain query string true "域名"
// @Param type query string true "域名类型 private action "
// @Router /manager/domain/create [post]
func (Domain) Create(c *gin.Context) {
	var req Request.DomainSave
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Logic.Domain{}.Create(req.Domain, req.TypeEd, "enable")
	Common2.ApiResponse{}.Success(c, "保存域名", gin.H{})
}

// @summary 域名-快速启用禁用域名
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query string true "域名id"
// @Param status query string true "状态 enable 启动 un_enable 禁用"
// @Router /manager/domain/enable_disable [post]
func (Domain) EnableDisable(c *gin.Context) {
	var req Request.DomainEnableDisable
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Logic.Domain{}.EnableDisable(req.Id, req.Status)
	Common2.ApiResponse{}.Success(c, "状态已跟新", gin.H{})
}

// @summary 域名-域名解绑
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param domain_id query string true "域名id"
// @Router /manager/domain/un_bind [post]
func (Domain) UnBind(c *gin.Context) {
	var req Request.DomainDelete
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	DomainCount := Logic.Domain{}.GetNoUsePrivateNum()
	if DomainCount == 0 {
		Common2.ApiResponse{}.Error(c, "可用域名不足，无法解绑", gin.H{})
		return
	}

	domain := Logic.Domain{}.Get(req.DomainId)
	update := map[string]interface{}{"bind_service_id": 0, "status": "down"}
	Base.MysqlConn.Model(&Common.Domain{}).Where("id = ?", req.DomainId).Updates(update)
	_ = Logic.Domain{}.Bind(domain.BindServiceId)

	Common2.ApiResponse{}.Success(c, "域名已解绑并下架", gin.H{})

	param := fmt.Sprintf("?service_id=%d&type=%s&content=%s", domain.BindServiceId, "ban", "域名拦截提醒")
	Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + param)
}
