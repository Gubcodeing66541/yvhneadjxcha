package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Model/Service"
	"server/Base"
)

type Black struct{}

// @summary 房间-拉黑指定用户
// @tags 客服后台
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Param type query string true "拉黑类型ip user"
// @Router /service/manager/black/add [post]
func (Black) Add(c *gin.Context) {
	var req Request.RoomBlack
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	if req.Type == "ip" {
		Logic.ServiceBlack{}.ServiceManagerIpBlack(roleId, req.Ip, req.Day, roleId)

	} else {
		Logic.ServiceBlack{}.ServiceManagerUserBlack(req.UserId, roleId, req.Day)
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

// @summary 房间-删除拉黑
// @tags 客服后台
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Param type query string true "拉黑类型ip user"
// @Router /service/manager/black/delete [post]
func (Black) Delete(c *gin.Context) {
	var req Request.RoomBlack
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	var serviceBlack Service.ServiceBlack
	Base.MysqlConn.Find(&serviceBlack, "id = ?", req.Id)

	//删除自己这条记录
	Base.MysqlConn.Delete(&Service.ServiceBlack{}, "id = ?", serviceBlack.Id)

	if serviceBlack.Type == "ip" {
		Base.MysqlConn.Delete(&Service.ServiceBlack{}, "service_manager_id = ? and ip = ? ", serviceBlack.ServiceManagerId, serviceBlack.Ip)
	}

	//其他的自動接觸拉黑
	var services []Service.Service
	Base.MysqlConn.Find(&services, "service_manager_id = ?", serviceBlack.ServiceManagerId)

	for _, v := range services {
		Base.MysqlConn.Model(&Service.ServiceRoom{}).Where("late_ip= ? and service_id =?", serviceBlack.Ip, v.ServiceId).Updates(gin.H{"is_black": 0})

	}

	Common.ApiResponse{}.Success(c, "删除成功", gin.H{})
}

func (Black) SearchUser(c *gin.Context) {

}
