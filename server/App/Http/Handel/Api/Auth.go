package Api

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Model/Service"
	"server/Base"
	"time"
)

type Auth struct{}

func (Auth) Login(c *gin.Context) {
	var req Request.RegisterAndLogin
	err := c.Bind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的账号密码", gin.H{})
		return
	}

	var serviceAuthModel Service.ServiceAuth
	Base.MysqlConn.Find(&serviceAuthModel, "username = ? and password = ?", req.Username, req.Password)
	if serviceAuthModel.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "账号密码有误", gin.H{})
		return
	}

	token := Common.Tools{}.EncodeToken(serviceAuthModel.ServiceId, "service", serviceAuthModel.ServiceId, 0)
	Common.ApiResponse{}.Success(c, "登录成功", gin.H{"token": token, "user_id": serviceAuthModel.ServiceId})
}

func (Auth) Register(c *gin.Context) {
	var req Request.RegisterAndLogin
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的账号密码", gin.H{})
		return
	}

	var serviceAuthModel Service.ServiceAuth
	Base.MysqlConn.Find(&serviceAuthModel, "username = ?", req.Username)
	if serviceAuthModel.ServiceId != 0 {
		Common.ApiResponse{}.Error(c, "账号已存在", gin.H{})
		return
	}

	Logic.Auth{}.Register(req.Username, req.Password, 0)
	Common.ApiResponse{}.Success(c, "注册成功", gin.H{})
}

func (Auth) UpdatePassword(c *gin.Context) {
	var req Request.UpdatePassword
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的账号密码", gin.H{})
		return
	}

	var serviceAuthModel Service.ServiceAuth
	Base.MysqlConn.Find(&serviceAuthModel, "username = ? and password = ?", req.Username, req.Password)
	if serviceAuthModel.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "账号密码有误", gin.H{})
		return
	}

	serviceAuthModel.Password = req.NewPassword
	serviceAuthModel.UpdateTime = time.Now()
	Base.MysqlConn.Save(&serviceAuthModel)
	Common.ApiResponse{}.Success(c, "密码修改成功", gin.H{})
}
