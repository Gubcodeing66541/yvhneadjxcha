package App

import (
	"github.com/gin-gonic/gin"
	"server/App/Http/Handel/Service"
	"server/App/Http/Handel/Servicemanager"
)

type ServiceManagerRoute struct{}

func (ServiceManagerRoute) BindRoute(s *gin.Engine) {
	s.POST("service/manager/auth/login", Servicemanager.Auth{}.Login)
	s.POST("service/manager/config", Service.Service{}.Config)
	service := s.Group("service/manager", ServiceMiddleWare())
	{
		service.POST("auth/reset_password", Servicemanager.Auth{}.ResetPassword)

		service.POST("data", Servicemanager.Manager{}.Data)
		service.POST("count", Servicemanager.Manager{}.Count)
		service.POST("count_room_detail", Servicemanager.Manager{}.CountRoomDetail)
		service.POST("count_message_detail", Servicemanager.Manager{}.CountMessageDetail)

		service.POST("info", Servicemanager.Manager{}.Info)
		service.POST("update", Servicemanager.Manager{}.Update)

		service.POST("message/add", Servicemanager.ServiceManagerMessage{}.Add)
		service.POST("message/list", Servicemanager.ServiceManagerMessage{}.List)
		service.POST("message/delete", Servicemanager.ServiceManagerMessage{}.Delete)
		service.POST("message/update", Servicemanager.ServiceManagerMessage{}.Update)

		service.POST("bot/info", Servicemanager.Bot{}.Info)
		service.POST("bot/update_info", Servicemanager.Bot{}.UpdateInfo)

		service.POST("bot/add", Servicemanager.Bot{}.Add)
		service.POST("bot/list", Servicemanager.Bot{}.List)
		service.POST("bot/delete", Servicemanager.Bot{}.Delete)
		service.POST("bot/update", Servicemanager.Bot{}.Update)

		service.POST("pay/recorder", Servicemanager.Pay{}.Recorder)
		service.POST("member/list", Servicemanager.Member{}.List)
		service.POST("member/create", Servicemanager.Member{}.Create)
		service.POST("member/create_list", Servicemanager.Member{}.CreateList)
		service.POST("member/update", Servicemanager.Member{}.Update)
		service.POST("member/delete", Servicemanager.Member{}.Delete)
		service.POST("member/renewal", Servicemanager.Member{}.Renewal)
		service.POST("member/renewal_all", Servicemanager.Member{}.RenewalAll)

		service.POST("users/list", Servicemanager.User{}.Users)
		service.POST("users/message", Servicemanager.User{}.Message)
		service.POST("users/black", Servicemanager.User{}.Black)

		service.POST("black/delete", Servicemanager.Black{}.Delete)
		service.POST("black/add", Servicemanager.Black{}.Add)
		service.POST("black/search_user", Servicemanager.Black{}.SearchUser)

	}
}
