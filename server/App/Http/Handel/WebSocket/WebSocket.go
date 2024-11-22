package WebSocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"server/App/Common"
)

type WebSocketConnect struct{}

// 升级websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @summary websocket连接
// @tags 公共
// @Param token header string true "认证token"
// @Router /api/websocket/conn [get]
func (WebSocketConnect) Conn(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Print("websocket-connect:err", err.Error())
		return
	}

	// 读取role_id 和 role_type
	roleId := Common.Tools{}.GetRoleId(c)
	roleType := Common.Tools{}.GetRoleType(c)
	GroupId := Common.Tools{}.GetRoleGroupId(c)

	// 客服信息
	serviceId := Common.Tools{}.GetServiceId(c)

	// 监控用户信息
	Event{}.OnConnect(roleType, roleId, GroupId, serviceId, ws)
}
