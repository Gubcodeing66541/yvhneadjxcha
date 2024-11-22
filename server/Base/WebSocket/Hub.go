package WebSocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"sync"
)

var Lock sync.RWMutex

type Hub struct {
	UserListMap        map[string]map[string]Connect
	UserConnGroupList  map[string]map[string]Connect
	UserConnIdGroupMap map[string]map[string]int
	ServiceBindUser    map[string]int
	ServiceBindGroup   map[string]int
}

// AddUser 添加用户
func (h *Hub) AddUser(c Connect) {
	// 如果没有映射表则先创建映射
	Lock.Lock()
	defer Lock.Unlock()
	if _, ok := h.UserListMap[c.UserId]; !ok {
		h.UserListMap[c.UserId] = map[string]Connect{}
	}

	// 如果没有绑定则创建新的并做数据绑定操作
	h.UserListMap[c.UserId][c.ConnId] = c

	// 用户组创建空数组
	h.UserConnIdGroupMap[c.ConnId] = map[string]int{}
}

// DelUser 删除用户
func (h *Hub) DelUser(c Connect) {

	// 在组中下线
	h.LeaveAllGroup(c, true)

	Lock.Lock()
	defer Lock.Unlock()

	// 如果没有查询到绑定的uid 也退出
	if _, ok := h.UserListMap[c.UserId][c.ConnId]; !ok {
		return
	}

	// 用户里面把自己下线
	delete(h.UserListMap[c.UserId], c.ConnId)
	delete(h.UserConnIdGroupMap, c.ConnId)
	delete(h.ServiceBindGroup, c.ConnId)

	// 如果所有用户都没有了则清空数组
	if len(h.UserListMap[c.UserId]) == 0 {
		delete(h.UserListMap, c.UserId)
		delete(h.ServiceBindUser, c.UserId)
		delete(h.ServiceBindGroup, c.UserId)
	}
}

// JoinGroupByUserId 通过userID加入组
func (h *Hub) JoinGroupByUserId(userId string, groupId string) {

	// 如果没有查询到组则创建新的组
	Lock.Lock()
	defer Lock.Unlock()

	if _, ok := h.UserListMap[userId]; !ok {
		return
	}

	for _, connect := range h.UserListMap[userId] {
		h.JoinGroup(connect, groupId, false)
	}
}

// LeaveAllGroupByUserId 退出所有组
func (h *Hub) LeaveAllGroupByUserId(userId string) {
	// 如果没有查询到组则创建新的组
	Lock.Lock()
	defer Lock.Unlock()

	if _, ok := h.UserListMap[userId]; !ok {
		return
	}

	for _, connect := range h.UserListMap[userId] {
		h.LeaveAllGroup(connect, false)
	}
}

// JoinGroup 加入组
func (h *Hub) JoinGroup(c Connect, groupId string, isLock bool) {
	if isLock {
		Lock.Lock()
		defer Lock.Unlock()
	}

	// 如果没有查询到组则创建新的组
	if _, ok := h.UserConnGroupList[groupId]; !ok {
		h.UserConnGroupList[groupId] = map[string]Connect{}
	}

	// 如果没有加组则加入组
	if _, ok := h.UserConnGroupList[groupId][c.ConnId]; !ok {
		h.UserConnGroupList[groupId][c.ConnId] = c
	}

	// 没有记录绑定的创建绑定关系
	if _, ok := h.UserConnIdGroupMap[c.ConnId]; !ok {
		h.UserConnIdGroupMap[c.ConnId] = map[string]int{}
	}

	// 没有加组记录则加组
	if _, ok := h.UserConnIdGroupMap[c.ConnId][groupId]; !ok {
		h.UserConnIdGroupMap[c.ConnId][groupId] = 1
	}
}

// LeaveAllGroup 退出所有组
func (h *Hub) LeaveAllGroup(c Connect, isLock bool) {
	if isLock {
		Lock.Lock()
		defer Lock.Unlock()
	}
	// 如果没有查询到绑定的uid 也退出
	if _, ok := h.UserConnIdGroupMap[c.ConnId]; !ok {
		return
	}

	// 循环退出所有的组
	for groupId := range h.UserConnIdGroupMap[c.ConnId] {
		// 查询到加租记录则删除组
		if _, ok := h.UserConnGroupList[groupId][c.ConnId]; ok {
			delete(h.UserConnGroupList[groupId], c.ConnId)
		}
	}

	// 清除自己绑定组的记录
	delete(h.UserConnIdGroupMap, c.ConnId)
}

// SendToUserId 指定用户发送消息
func (h *Hub) SendToUserId(userId string, message []byte) {
	Lock.Lock()
	defer Lock.Unlock()
	// 不在线则跳过发送
	if _, ok := h.UserListMap[userId]; !ok {
		return
	}

	// 给在线的userID 推送消息
	var userConn Connect
	for _, userConn = range h.UserListMap[userId] {
		err := userConn.Conn.WriteMessage(1, message)
		if err != nil {
			fmt.Print("websocket-err: SendToUserId", err.Error())
		}
	}
}

// SendToConnId 指定用户发送消息
func (h *Hub) SendToConnId(userId string, connId string, message []byte) {
	Lock.Lock()
	defer Lock.Unlock()
	// 不在线则跳过发送
	if _, ok := h.UserListMap[userId][connId]; !ok {
		return
	}

	err := h.UserListMap[userId][connId].Conn.WriteMessage(1, message)
	if err != nil {
		fmt.Print("websocket-err: SendToConnId", err.Error())
	}
	return
}

// SendToGroupId 指定组发送消息
func (h *Hub) SendToGroupId(groupId string, message []byte) {
	Lock.Lock()
	defer Lock.Unlock()
	// 如果没有查询到组则没有推送
	if _, ok := h.UserConnGroupList[groupId]; !ok {
		return
	}

	// 循环推送消息
	var groupInfo Connect
	for _, groupInfo = range h.UserConnGroupList[groupId] {
		err := groupInfo.Conn.WriteMessage(1, message)
		if err != nil {
			fmt.Print("websocket-err: SendToGroupId", err.Error())
		}
	}
}

// UserIdIsOnline 是否在线 true 在线 false 不在线
func (h *Hub) UserIdIsOnline(userId string) int {
	Lock.RLock()
	defer Lock.RUnlock()
	_, ok := h.UserListMap[userId]
	if ok {
		return 1
	}
	return 0
}

// GetOnlineCount GetAllConn 获取在线用户数
func (h *Hub) GetOnlineCount() int {
	Lock.RLock()
	defer Lock.RUnlock()
	return len(h.UserConnIdGroupMap)
}

// GetAllStatus GetAllConn 获取在线用户数
func (h *Hub) GetAllStatus() map[string]interface{} {
	Lock.RLock()
	defer Lock.RUnlock()

	return map[string]interface{}{"h": gin.H{
		"UserConnIdGroupMap": h.UserConnIdGroupMap,
		"ServiceBindUser":    h.ServiceBindUser,
		"ServiceBindGroup":   h.ServiceBindGroup,
	}}
}

// GetAllConn 获取所有连接
func (h *Hub) GetAllConn() ([]int, []int) {
	Lock.RLock()
	defer Lock.RUnlock()

	user := make([]int, 0)
	service := make([]int, 0)
	for socketUserId, _ := range h.UserListMap {
		temp := strings.Split(socketUserId, ":")
		RoleId, err := strconv.Atoi(temp[1])
		if err != nil {
			continue
		}

		if temp[0] == "service" {
			service = append(service, RoleId)

		}

		if temp[0] == "user" {
			user = append(user, RoleId)
		}
	}
	return service, user

	//return h.UserListMap
}

func (h *Hub) BindUser(ServiceUserId string, UserId int) {
	Lock.Lock()
	defer Lock.Unlock()
	h.ServiceBindUser[ServiceUserId] = UserId
}

func (h *Hub) BindGroup(ServiceId string, groupId int) {
	Lock.Lock()
	defer Lock.Unlock()
	h.ServiceBindGroup[ServiceId] = groupId
}

func (h *Hub) GetBindUser(ServiceId string) int {
	Lock.RLock()
	defer Lock.RUnlock()
	userId, ok := h.ServiceBindUser[ServiceId]
	if ok {
		return userId
	}
	return 0
}

func (h *Hub) GetBindGroup(ServiceId string) int {
	Lock.RLock()
	defer Lock.RUnlock()
	groupId, ok := h.ServiceBindGroup[ServiceId]
	if ok {
		return groupId
	}
	return 0
}

// Run 启动服务
func (h *Hub) Run() {

}
