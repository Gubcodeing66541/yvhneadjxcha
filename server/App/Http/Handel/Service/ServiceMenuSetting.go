package Service

import (
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Request"
	Service3 "server/App/Http/Request/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
)

type ServiceMenuSetting struct{}

// @summary 菜单管理-创建
// @tags 客服系统
// @Param token header string true "认证token"
// @Param title query string true "标题"
// @Param content query string true "内容"
// @Param action query string true "操作"
// @Param tag query string true "提示"
// @Param sort query string true "排序"
// @Router /service/menu_setting/create [post]
func (ServiceMenuSetting) Create(c *gin.Context) {
	var req Service3.ServiceMenuSetting
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{"req": req})
		return
	}

	serviceId := Common2.Tools{}.GetServiceId(c)
	Base.MysqlConn.Create(&Service2.ServiceMenuSetting{
		Tag:       req.Tag,
		ServiceId: serviceId, Title: req.Title, Content: req.Content, Action: req.Action, Sort: req.Sort,
	})

	Common2.ApiResponse{}.Success(c, "创建成功", gin.H{})
}

// @summary 菜单管理-删除消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param id query int true "删除的ID"
// @Router /service/menu_setting/delete [post]
func (ServiceMenuSetting) Delete(c *gin.Context) {
	var req Request.DeleteServiceMessage
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)
	Base.MysqlConn.Delete(&Service2.ServiceMenuSetting{}, "id = ? and service_id = ? ", req.Id, ServiceId)
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{"req": req, "test": 11111})
}

// @summary 菜单管理-修改消息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param id query string true "ID"
// @Param title query string true "标题"
// @Param content query string true "内容"
// @Param action query string true "操作"
// @Param tag query string true "提示"
// @Param sort query string true "排序"
// @Router /service/menu_setting/update [post]
func (ServiceMenuSetting) Update(c *gin.Context) {
	var req Service3.ServiceMenuSetting
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误了", gin.H{"req": req, "r": 1})
		return
	}
	ServiceId := Common2.Tools{}.GetServiceId(c)
	Base.MysqlConn.Model(&Service2.ServiceMenuSetting{}).Where("id = ? and service_id = ?", req.Id, ServiceId).
		Updates(Service2.ServiceMenuSetting{
			Title: req.Title, Content: req.Content, Action: req.Action, Sort: req.Sort, Tag: req.Tag,
		})

	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 菜单管理-消息列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/menu_setting/list [post]
func (ServiceMenuSetting) List(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	ServiceId := Common2.Tools{}.GetServiceId(c)
	tel := Base.MysqlConn.Model(&Service2.ServiceMenuSetting{}).Where("service_id = ?", ServiceId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Service2.ServiceMenuSetting
	tel.Order("sort asc").Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Find(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}
