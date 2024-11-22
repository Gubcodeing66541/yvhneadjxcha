package Manage

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Request"
	ServiceManager2 "server/App/Model/ServiceManager"
	"server/Base"
)

type Auth struct{}

// @summary 登录账号
// @tags 客服系统总后台
// @description 测试账号 admin 密码 123456`
// @Param username query string true "账号"
// @Param password query string true "密码"
// @Router /manager/auth/login [post]
func (Auth) Login(c *gin.Context) {
	var req Request.RegisterAndLogin
	err := c.Bind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的账号密码", gin.H{})
		return
	}
	var serviceManagerAuth ServiceManager2.ServiceManagerAuth

	Base.MysqlConn.Find(&serviceManagerAuth, "username=? and password = ?", req.Username, req.Password)
	if serviceManagerAuth.ServiceManagerId == 0 {
		Common.ApiResponse{}.Error(c, "账号或密码错误", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(1, "manage", 0, 0)
	Common.ApiResponse{}.Success(c, "登录成功", gin.H{"token": token})
}

// @summary 修改密码
// @tags 客服系统总后台
// @description 原密码 admin 新密码 123456`
// @Param password query string true "密码"
// @Param new_password query string true "新密码"
// @Router /manager/auth/login [post]
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
	var serviceManagerModel ServiceManager2.ServiceManagerAuth
	Base.MysqlConn.Find(&serviceManagerModel, "service_manager_id = ? and password = ? ", roleId, req.Password)
	if serviceManagerModel.ServiceManagerId == 0 {
		Common.ApiResponse{}.Error(c, "账号或密码错误,请检查", gin.H{})
		return
	}

	Base.MysqlConn.Model(&ServiceManager2.ServiceManagerAuth{}).Where("service_manager_id = ?", roleId).Updates(gin.H{"password": req.NewPassword})
	Common.ApiResponse{}.Success(c, "密码修改成功", gin.H{})

}
