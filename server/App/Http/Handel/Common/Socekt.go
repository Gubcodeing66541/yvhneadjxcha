package Common

import (
	"github.com/gin-gonic/gin"
	Common2 "server/App/Common"
	"server/App/Http/Request"
	"server/App/Model/Service"
	"server/Base"
)

type Socket struct{}

// GetAllByManager 获取所有的连接数据
func (Socket) GetAllByManager(c *gin.Context) {
	serviceIds, user := Common2.Socket{}.GetAll()
	serviceCount, userCount := len(serviceIds), len(user)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{
		"service_ids":   serviceIds,
		"user_ids":      user,
		"service_count": serviceCount,
		"user_count":    userCount,
		"socket_count":  Base.WebsocketHub.GetOnlineCount(),
	})
}

// GetAllByManager 获取所有的连接数据
func (Socket) GetAllByServiceManager(c *gin.Context) {
	serviceIds, user := Common2.Socket{}.GetAll()
	serviceCount, userCount := len(serviceIds), len(user)

	var services []Service.Service
	Base.MysqlConn.Where("service_id in (?)", serviceIds).Find(&services)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{
		"service_ids":    serviceIds,
		"service_count":  serviceCount,
		"user_count":     userCount,
		"socket_count":   Base.WebsocketHub.GetOnlineCount(),
		"online_service": services,
	})
}

// SendToServiceSocket 获取所有的连接数据
func (Socket) SendToServiceSocket(c *gin.Context) {
	var req Request.SocketMsg
	err := c.ShouldBind(&req)
	if err != nil {
		return
	}
	Common2.ApiResponse{}.SendMsgToService(req.ServiceId, req.Type, req.Content)
	Common2.ApiResponse{}.Success(c, "发送成功", gin.H{"req": req})
}
