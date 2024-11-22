package App

import (
	"github.com/gin-gonic/gin"
	"server/App/Http/Handel/Manage"
	"server/App/Http/Handel/Service"
)

type ManagerRoute struct{}

func (ManagerRoute) BindRoute(s *gin.Engine) {
	s.POST("manager/auth/login", Manage.Auth{}.Login)

	s.POST("manager/config", Service.Service{}.Config)
	manager := s.Group("manager", ManagerMiddleWare())
	{
		manager.POST("auth/reset_password", Manage.Auth{}.ResetPassword)
		manager.POST("count", Manage.Service{}.Count)

		manager.POST("count_service_list", Manage.Service{}.CountServiceList)

		manager.POST("pay/recorder", Manage.Pay{}.Recorder)

		//客服管理
		manager.POST("service_manager/list", Manage.ServiceManager{}.List)
		manager.POST("service_manager/create", Manage.ServiceManager{}.Create)
		manager.POST("service_manager/renew", Manage.ServiceManager{}.ReNew)
		manager.POST("service_manager/get_service_list", Manage.ServiceManager{}.GetServiceList)
		manager.POST("service_manager/delete", Manage.ServiceManager{}.Delete)
		manager.POST("service_manager/reset_password", Manage.ServiceManager{}.ResetPassword)
		manager.POST("service_manager/ban", Manage.ServiceManager{}.Ban)

		//客服管理
		manager.POST("service/list", Manage.Service{}.ServiceList)
		manager.POST("service/create", Manage.Service{}.ServiceCreate)
		manager.POST("service/bach_create", Manage.Service{}.ServiceBachCreate)
		manager.POST("service/renewal", Manage.Service{}.Renewal)

		//域名管理
		manager.POST("domain/list", Manage.Domain{}.List)
		manager.POST("domain/query_by_id", Manage.Domain{}.QueryById)
		manager.POST("domain/delete", Manage.Domain{}.Delete)
		manager.POST("domain/update", Manage.Domain{}.Update)
		manager.POST("domain/create", Manage.Domain{}.Create)
		manager.POST("domain/enable_disable", Manage.Domain{}.EnableDisable)
		manager.POST("domain/un_bind", Manage.Domain{}.UnBind)

		//公众号管理
		manager.POST("wechat_auths/list", Manage.WeChatAuths{}.List)
		manager.POST("wechat_auths/enable_disable", Manage.WeChatAuths{}.OpenOrClose)
		manager.POST("wechat_auths/delete", Manage.WeChatAuths{}.Delete)
		manager.POST("wechat_auths/update", Manage.WeChatAuths{}.Update)
		manager.POST("wechat_auths/create", Manage.WeChatAuths{}.Create)
		manager.POST("wechat_auths/switch", Manage.WeChatAuths{}.Switch)

		manager.POST("service/get_order", Manage.Service{}.GetServiceOrder)
		manager.POST("service/get_order_info", Manage.Service{}.GetServiceOrderInfo)
		manager.POST("service/bind_domain", Manage.Service{}.BindDomain)
		manager.POST("service/change_bind_domain", Manage.Service{}.ChangeBindDomain)
		manager.POST("service/get_service_domain", Manage.Service{}.GetServiceDomain)

		manager.POST("config/get", Manage.Config{}.Get)
		manager.POST("config/update", Manage.Config{}.Update)

		manager.POST("setting/list", Manage.Setting{}.List)
		manager.POST("setting/create", Manage.Setting{}.Create)
		manager.POST("setting/update", Manage.Setting{}.Update)
		manager.POST("setting/delete", Manage.Setting{}.Delete)

		manager.POST("users/list", Manage.User{}.List)

		manager.POST("clear/message", Manage.User{}.ClearMessage) // 清理聊天记录
		
	}
}
