package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Model/ServiceManager"
	"server/Base"
)

type Auth struct{}

// @summary 基础-登录
// @tags 客服后台
// @Description 服务初始连接测试 --> 接口描述
// @Param username query string true "登录账号" default(1)
// @Param password query string true "登录密码" default(1)
// @Router /service/manager/auth/login [post]
func (Auth) Login(c *gin.Context) {
	var req Request.RegisterAndLogin
	err := c.ShouldBind(&req)
	if err != nil || (req.Username == "" || req.Password == "") {
		Common.ApiResponse{}.Error(c, "请输入正确的账号密码", gin.H{})
		return
	}

	//查询账号密码
	var serviceManagerModel ServiceManager.ServiceManager
	Base.MysqlConn.Find(&serviceManagerModel, "member = ? and password = ? ", req.Username, req.Password)
	if serviceManagerModel.ServiceManagerId == 0 {
		Common.ApiResponse{}.Error(c, "账号或密码错误,请检查后在登录", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(serviceManagerModel.ServiceManagerId, "service_manager", 0, 0)

	//更新IP
	Base.MysqlConn.Model(&ServiceManager.ServiceManager{}).Where("service_manager_id =?", serviceManagerModel.ServiceManagerId).Updates(ServiceManager.ServiceManager{
		Ip: c.ClientIP(),
	})

	Common.ApiResponse{}.Success(c, "登录成功", gin.H{"token": token})
}

// @summary 基础-修改密码
// @tags 客服后台
// @Param token header string true "认证token"
// @Description 服务初始连接测试 --> 接口描述
// @Param username query string true "登录账号" default(1)
// @Param password query string true "登录密码" default(1)
// @Param new_password query string true "登录密码" default(1)
// @Router /service/manager/auth/reset_password [post]
func (Auth) ResetPassword(c *gin.Context) {
	var req Request.UpdatePassword
	err := c.ShouldBind(&req)
	if err != nil || (req.NewPassword == "" || req.Password == "") {
		Common.ApiResponse{}.Error(c, "请输入全", gin.H{})
		return
	}

	//token获取账号
	roleId := Common.Tools{}.GetRoleId(c)

	//查询账号密码
	var serviceManagerModel ServiceManager.ServiceManager
	Base.MysqlConn.Find(&serviceManagerModel, "service_manager_id = ? and password = ? ", roleId, req.Password)
	if serviceManagerModel.ServiceManagerId == 0 {
		Common.ApiResponse{}.Error(c, "账号或密码错误,请检查", gin.H{})
		return
	}

	Base.MysqlConn.Model(&ServiceManager.ServiceManager{}).Where("service_manager_id = ?", roleId).Updates(gin.H{"password": req.NewPassword})
	Common.ApiResponse{}.Success(c, "密码修改成功", gin.H{})
}
