package App

import (
	"fmt"
	"net/http"
	"server/App/Common"
	"server/App/Http/Constant"

	"github.com/gin-gonic/gin"
)

type HttpRoute struct{}

func (HttpRoute) BindRoute(s *gin.Engine) {

	s.Use(Cors())
	ApiRoute{}.BindRoute(s)
	ServiceRoute{}.BindRoute(s)
	ServiceManagerRoute{}.BindRoute(s)
	UserRoute{}.BindRoute(s)
	ManagerRoute{}.BindRoute(s)
}

func ApiMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.GetHeader("token")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		if userInfo.RoleType == "" {
			Common.ApiResponse{}.NoAuth(c, gin.H{"role_type": userInfo.RoleType})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}

func UserMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.GetHeader("token")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		if userInfo.RoleType != "user" {
			Common.ApiResponse{}.NoAuth(c, gin.H{"role_type": userInfo.RoleType})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}

func ManagerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// token
		token := c.GetHeader("token")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		if userInfo.RoleType != "manage" {
			Common.ApiResponse{}.NoAuth(c, gin.H{"role_type": userInfo.RoleType})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}

func ServiceMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.RoleId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}

func ServiceManagerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error()})
			fmt.Println("权限不足")
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.RoleId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Next()
	}
}

func WebSocketMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		var userInfo Constant.UserAuthToken
		err := Common.Tools{}.DecodeToken(token, &userInfo)
		if err != nil {
			Common.ApiResponse{}.NoAuth(c, gin.H{"err": err.Error(), "token": token, "type": "websocket"})
			c.Abort()
			return
		}

		c.Set("service_id", userInfo.ServiceId)
		c.Set("role_id", userInfo.RoleId)
		c.Set("role_type", userInfo.RoleType)
		c.Set("group_id", userInfo.GroupId)
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Requested-With,XMLHttpRequest")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Token ,token")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusOK)
		}
		context.Next()
	}
}
