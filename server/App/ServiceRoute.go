package App

import (
	"github.com/gin-gonic/gin"
	"server/App/Http/Handel/Common"
	"server/App/Http/Handel/Service"
)

type ServiceRoute struct{}

func (ServiceRoute) BindRoute(s *gin.Engine) {
	s.POST("service/auth/login", Service.Service{}.Login)

	s.POST("service/common/count", Common.Common{}.Count)
	s.POST("service/common/api/oss_config", Common.Common{}.Oss)

	s.POST("service/config", Service.Service{}.Config)

	service := s.Group("service", ServiceMiddleWare())

	{
		service.POST("reset_qrcode", Service.Service{}.ResetQrcode)
		service.POST("update_qrcode", Service.Service{}.UpdateQrcode)

		service.POST("service_manager/message/list", Service.Service{}.ServiceManagerMessageList)

		service.POST("users", Service.Service{}.UserList)
		service.POST("count", Service.Service{}.Count)

		service.POST("menu_setting/create", Service.ServiceMenuSetting{}.Create)
		service.POST("menu_setting/update", Service.ServiceMenuSetting{}.Update)
		service.POST("menu_setting/list", Service.ServiceMenuSetting{}.List)
		service.POST("menu_setting/delete", Service.ServiceMenuSetting{}.Delete)

		service.POST("notice_setting/update", Service.ServiceNoticeSetting{}.Update)
		service.POST("notice_setting/info", Service.ServiceNoticeSetting{}.Info)

		service.POST("setting", Service.Service{}.SettingList)
		service.POST("info", Service.Service{}.Info)
		service.POST("update", Service.Service{}.Update)
		service.POST("/del", Service.Service{}.DelService) //删除用户

		service.POST("message/send_all", Service.Message{}.SendAll)
		service.POST("message/send_to_user", Service.Message{}.SendToUser)
		service.POST("message/list", Service.Message{}.List)
		service.POST("message/update", Service.Message{}.Update)
		service.POST("message/remove_msg", Service.Message{}.RemoveMessage)
		service.POST("message/clear_message", Service.Message{}.ClearMessage)

		service.POST("rooms/list", Service.ServiceRooms{}.List)
		service.POST("rooms/detail", Service.ServiceRooms{}.Detail)
		service.POST("rooms/update", Service.ServiceRooms{}.Update)
		service.POST("rooms/top", Service.ServiceRooms{}.Top)
		service.POST("rooms/black", Service.ServiceRooms{}.Black)
		service.POST("rooms/black_list", Service.ServiceRooms{}.BlackList)
		service.POST("rooms/count", Service.ServiceRooms{}.Count)
		service.POST("rooms/rename", Service.ServiceRooms{}.Rename) //修改昵称
		service.POST("rooms/delete_day", Service.ServiceRooms{}.DeleteDay)
		service.POST("rooms/end", Service.ServiceRooms{}.End) //修改昵称
	}

	serviceMessage := s.Group("service/service_message", ServiceMiddleWare())
	{
		serviceMessage.POST("create", Service.ServiceMessage{}.Create)
		serviceMessage.POST("delete", Service.ServiceMessage{}.Delete)
		serviceMessage.POST("update", Service.ServiceMessage{}.Update)
		serviceMessage.POST("list", Service.ServiceMessage{}.List)
		serviceMessage.POST("get", Service.ServiceMessage{}.GetById)
		serviceMessage.POST("swap", Service.ServiceMessage{}.Swap)
	}
}
