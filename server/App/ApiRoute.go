package App

import (
	"net/http"
	Common2 "server/App/Common"
	"server/App/Constant"
	"server/App/Http/Handel"
	"server/App/Http/Handel/Api"
	"server/App/Http/Handel/Common"
	"server/App/Http/Handel/WebSocket"
	"server/Base"

	"github.com/gin-gonic/gin"
)

type ApiRoute struct{}

func (ApiRoute) BindRoute(s *gin.Engine) {
	s.StaticFS("service", http.Dir("./Tel/dist/service"))
	s.StaticFS("service_manage", http.Dir("./Tel/dist/service_manage"))
	s.StaticFS("manage", http.Dir("./Tel/dist/manage"))

	s.StaticFS("users", http.Dir("./Tel/dist/user"))
	s.StaticFS("common", http.Dir("./Tel/common"))
	s.StaticFS("/s", http.Dir("./static"))

	//
	s.LoadHTMLGlob("Tel/dist/**/*.html")
	s.LoadHTMLGlob("Test/build/*.html")

	//
	s.StaticFS("/static", http.Dir("./static"))
	s.StaticFS("/head", http.Dir("./static/head"))
	//s.StaticFS("/upload", http.Dir("./static/upload"))

	s.GET(":filename", Common.Common{}.WeChatFile)

	s.GET("test/create", func(context *gin.Context) {
		Common2.CreateDomain()
		context.String(http.StatusOK, "ok")
	})

	common := s.Group("common")
	{
		common.POST("/oss/qiniu_token", Common.Qiniu{}.QiniuToken)
	}

	// 解决跨域
	api := s.Group("api")
	{
		api.GET("socket/send_to_service_socket", Common.Socket{}.SendToServiceSocket)

		api.POST("auth/login", Api.Auth{}.Login)
		api.POST("auth/register", Api.Auth{}.Register)
		api.POST("auth/up_password", Api.Auth{}.UpdatePassword)
		api.GET("websocket/conn", WebSocketMiddleWare(), WebSocket.WebSocketConnect{}.Conn)
		system := api.Group("system")
		{
			system.GET("status_all", Common.Socket{}.GetAllByManager)
			system.GET("status", Handel.System{}.Status)
			system.POST("clear_cache", Handel.System{}.ClearCache)
			system.POST("upload", ApiMiddleWare(), Handel.System{}.Upload)
			system.POST("upload_image", ApiMiddleWare(), Handel.System{}.UploadImage)
			system.GET("action", Handel.System{}.Action)
			system.GET("location", Handel.System{}.Location)
		}
	}
	api.GET("test/push", func(context *gin.Context) {
		err := Base.Producer.Publish(Constant.Topic, []byte("message"))
		if err != nil {
			println(err.Error())
			context.String(http.StatusOK, "err")
			return
		}
		context.String(http.StatusOK, "ok")
	})

	api.POST("tools/copy", Api.Tools{}.Copy)
	api.POST("tools/search", Api.Tools{}.Search)
	api.POST("tools/service_count", Api.Tools{}.ServiceCount)
	api.POST("tools/getcode", Api.Tools{}.GetCode)
	api.POST("tools/fix_domain", Api.Tools{}.FixDomain)
	api.POST("tools/reset_domain", Api.Tools{}.ResetDomain)
	api.POST("tools/count", Api.Tools{}.Count)

}
