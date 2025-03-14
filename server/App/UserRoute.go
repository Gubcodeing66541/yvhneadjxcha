package App

import (
	"server/App/Http/Handel/Common"
	"server/App/Http/Handel/User"

	"github.com/gin-gonic/gin"
)

type UserRoute struct{}

func (UserRoute) BindRoute(s *gin.Engine) {
	user := s.Group("user")
	{
		user.GET("code/actions", User.LocalAuth{}.CodeActions) // 新二维码图片仅仅只是展示图片信息
		user.GET("code/img", User.LocalAuth{}.Img)             // 图片显示

		user.POST("common/api/oss_config", Common.Common{}.Oss)

		//微信公众号逻辑
		user.GET("auth/web/:code", User.Auth{}.Web)
		user.GET("auth/join", User.Auth{}.Join)
		user.GET("auth/transfer", User.Auth{}.Transfer)
		user.GET("auth/we_chat_auth", User.Auth{}.WeChatAuth)

		//uuid逻辑
		user.GET("auth/test/:code/:uuid", User.Auth{}.Test) //test

		user.GET("auth/no_wechat_auth", User.Auth{}.NoWeChatAuth)                     //入口
		user.GET("auth/no_wechat_transfer/:code/:uuid", User.Auth{}.NoWechatTransfer) //中转
		user.GET("auth/no_wechat_action/:code/:uuid", User.Auth{}.NoWechatAction)     //落地

		user.GET("auth/local_storage/join", User.LocalAuth{}.Join)                                  //入口
		user.GET("auth/local_storage/transfer_action/:code/:uuid", User.LocalAuth{}.TransferAction) //中转
		user.GET("auth/local_storage/transfer/:code/:uuid", User.LocalAuth{}.Transfer)              //中转
		user.GET("auth/local_storage/action/:code/:uuid", User.LocalAuth{}.Action)                  //落地
		user.GET("auth/local_storage/show/:code/:uuid", User.LocalAuth{}.Show)                      //展示
		user.GET("auth/local_storage/location", User.LocalAuth{}.Location)                          //入口

		user.POST("auth/local_storage/bind_uuid/:new_uuid/:current_uuid", User.LocalAuth{}.BindUUid) //中转

		user.GET("auth/error/:msg", User.Auth{}.Error)

		user.GET("auth/action/:token", User.Auth{}.Action)
		user.POST("auth/token", User.Auth{}.Token)
		user.POST("auth/info", UserMiddleWare(), User.Auth{}.Info)
		user.POST("auth/send", UserMiddleWare(), User.Auth{}.Send)
		user.POST("message/list", UserMiddleWare(), User.Auth{}.List)
		user.POST("message/button", UserMiddleWare(), User.Auth{}.Button)

		// 入口和落地
		user.GET("auth/local_storage/join_new", User.LocalAuth{}.JoinNew) //入口

		user.POST("oauth/action", User.OtherAuth.Action)                   // 新落地 user/oauth/action
		user.POST("oauth/domain", UserMiddleWare(), User.OtherAuth.Domain) // 新落地 user/oauth/domain
		user.POST("oauth/token", User.OtherAuth.Token)                     // 新落地 uuid换token

		user.GET("oauth/show_join", User.OtherAuth.ShowJoin)     // 新落地join
		user.GET("oauth/show_action", User.OtherAuth.ShowAction) // 新落地action
		user.GET("j", User.OtherAuth.ShowJoin)                   // 新落地join

	}
}
