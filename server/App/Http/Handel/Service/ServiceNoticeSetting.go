package Service

import (
	"github.com/gin-gonic/gin"
	Common2 "server/App/Common"
	Service3 "server/App/Http/Request/Service"
	Service2 "server/App/Model/Service"
	"server/Base"
)

type ServiceNoticeSetting struct{}

// @summary 公告-修改
// @tags 客服系统
// @Param token header string true "认证token"
// @Param is_show query string true "启用状态"
// @Param image query string true "图片"
// @Param text query string true "内容"
// @Router /service/notice_setting/update [post]
func (ServiceNoticeSetting) Update(c *gin.Context) {
	var req Service3.ServiceNoticeSetting
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	serviceId := Common2.Tools{}.GetRoleId(c)

	findServiceNoticeSetting := Service2.ServiceNoticeSetting{}
	Base.MysqlConn.Where(" service_id = ?", serviceId).Find(&findServiceNoticeSetting)
	if findServiceNoticeSetting.Id == 0 {
		findServiceNoticeSetting.Text = req.Text
		findServiceNoticeSetting.Image = req.Image
		findServiceNoticeSetting.IsShow = req.IsShow
		findServiceNoticeSetting.ServiceId = serviceId
		Base.MysqlConn.Save(&findServiceNoticeSetting)
		//return
	}

	Base.MysqlConn.Model(&Service2.ServiceNoticeSetting{}).Where(" service_id = ?", serviceId).Updates(
		gin.H{"is_show": req.IsShow, "image": req.Image, "text": req.Text})
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 公告-详细信息
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/notice_setting/info [post]
func (ServiceNoticeSetting) Info(c *gin.Context) {
	ServiceId := Common2.Tools{}.GetServiceId(c)
	var notice Service3.ServiceNoticeSetting
	Base.MysqlConn.Find(&notice, "service_id = ?", ServiceId)
	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{"notice": notice})
}
