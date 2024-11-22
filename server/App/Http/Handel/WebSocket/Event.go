package WebSocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"server/App/Common"
	"server/Base"
	WebSocket2 "server/Base/WebSocket"
	"time"
)

type Event struct{}

func (Event) OnConnect(roleType string, roleId int, groupId string, serviceId int, ws *websocket.Conn) {

	username := fmt.Sprintf("%s:%d", roleType, roleId)

	// 生成ConnId
	connId := fmt.Sprintf("user_id_%s_time_%s_rand_number_%d", username, time.Now(), rand.Intn(999999))
	connId = Common.Tools{}.Md5(connId)

	//Binary 1 消息采用二进制格式。  Close 2 因为收到关闭的消息，接受已完成。 Text 0 该消息是明文形式。
	conn := WebSocket2.Connect{UserId: username, Conn: ws, ConnId: connId}
	Base.WebsocketHub.AddUser(conn)

	if groupId != "" {
		Base.WebsocketHub.JoinGroup(conn, groupId, true)
	}

	// 推送上线状态和离线状态
	if roleType == "user" {
		defer Common.ApiResponse{}.SendMsgToService(serviceId, "leave", gin.H{"user_id": roleId})
		Common.ApiResponse{}.SendMsgToService(serviceId, "online", gin.H{"user_id": roleId})
	} else {
		Base.WebsocketHub.BindGroup(connId, 0)

	}

	//返回前关闭 和删除用户连接数据
	defer ws.Close()
	defer Base.WebsocketHub.DelUser(conn)
	defer Event{}.OnClose(conn)

	//读取ws中的数据 转交给业务处理
	for {
		// 设置超时时间
		err := ws.SetReadDeadline(time.Now().Add(20 * time.Second))
		if err != nil {
			fmt.Println("超时检测错误")
			return
		}
		//
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("socket连接已断开", err, err.Error())
			break
		}

		Event{}.OnMessage(conn, message)
	}
}

func (Event) OnMessage(conn WebSocket2.Connect, message []byte) {
}

func (Event) OnClose(conn WebSocket2.Connect) {

}
