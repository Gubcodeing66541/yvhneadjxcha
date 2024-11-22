package Logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Message2 "server/App/Model/Message"
	Service2 "server/App/Model/Service"
	User2 "server/App/Model/User"
	"server/Base"
	"time"
)

type ServiceRoom struct {
}

func (ServiceRoom) Get(user User2.User, serviceId int, Ip string, Drive string) Service2.ServiceRoom {
	var serverRoom Service2.ServiceRoom
	RoomId := Common.Tools{}.ConvertUserMessageRoomId(user.UserId, serviceId)
	Base.MysqlConn.Find(&serverRoom, "service_id = ? and user_id = ?", serviceId, user.UserId)

	if serverRoom.Id == 0 {
		now := time.Now()
		var hello []Service2.ServiceMessage
		Base.MysqlConn.Find(&hello, "service_id = ? and type = ? and status = 'enable'", serviceId, "hello")
		LateMsg, LateType, lateId := "", "", 0
		for _, item := range hello {
			model := &Message2.Message{
				RoomId: RoomId, From: serviceId, To: user.UserId, Type: item.MsgType, Content: item.MsgInfo,
				SendRole: "hello", CreateTime: now, IsRead: 1, UserId: user.UserId, ServiceId: serviceId}
			Base.MysqlConn.Create(&model)
			LateType, LateMsg = item.MsgType, item.MsgInfo
			lateId = model.Id
		}
		serverRoom = Service2.ServiceRoom{
			RoomId: RoomId, ServiceId: serviceId, LateUserReadId: lateId, UserId: user.UserId, LateType: LateType,
			LateMsg: LateMsg, CreateTime: now, UpdateTime: now, LateId: lateId}
		Base.MysqlConn.Create(&serverRoom)
		Base.MysqlConn.Create(&Service2.ServiceRoomDetail{ServiceId: serviceId, UserId: user.UserId, RoomId: RoomId, IP: Ip, CreateTime: now, Drive: Drive})
	}

	go func(Ip string) {
		userAddr, _ := Common.Tools{}.IPInfo(Ip)
		fmt.Println("------------------------------用户的获取到的异步的地址是", RoomId, "---", userAddr.Addr)
		Base.MysqlConn.Model(&Service2.ServiceRoomDetail{}).
			Where("service_id = ? and user_id = ?", serviceId, user.UserId).Updates(gin.H{"map": userAddr.Addr})

	}(Ip)
	return serverRoom
}

func (ServiceRoom) List(serviceId int, req Request.ServiceRoomList, isRoomsSearch bool) gin.H {
	var model []Response.ServiceRoom

	where := fmt.Sprintf("service_rooms.service_id = %d ", serviceId)

	if isRoomsSearch {
		where += fmt.Sprintf(" and is_delete = 0 ")
	}

	if req.UserName != "" {
		where += fmt.Sprintf(" and (users.user_name like '%%%s%%' or service_rooms.rename like '%%%s%%')", req.UserName, req.UserName)
	}

	if req.Type == "server_read" {
		where += fmt.Sprintf(" and service_rooms.late_role = 'service' ")
	}

	if req.Type == "server_no_read" {
		where += fmt.Sprintf(" and service_rooms.late_role = 'user' ")
	}

	if req.Type == "server_no_read_count" {
		where += fmt.Sprintf(" and service_rooms.service_no_read > 0 ")
	}

	if req.Type == "top" {
		where += fmt.Sprintf(" and service_rooms.is_top = 1 ")
	}

	if req.Type == "black" {
		where += fmt.Sprintf(" and service_rooms.is_black = %d ", 1)
	} else {
		where += fmt.Sprintf(" and service_rooms.is_black = %d ", 0)
	}

	show := "service_rooms.id,service_rooms.late_type,service_rooms.is_top,service_rooms.user_no_read,service_rooms.late_msg,service_rooms.service_no_read,service_rooms.update_time,service_rooms.rename,"
	show += "users.user_id,users.user_name,users.user_head"
	join := "join users on service_rooms.user_id = users.user_id"

	tel := Base.MysqlConn.Table("service_rooms").Select(show).Where(where).Joins(join)

	if isRoomsSearch {
		tel = tel.Order("service_rooms.is_top desc,service_rooms.update_time desc")
	} else {
		tel = tel.Order("id desc")
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel.Offset((req.Page - 1) * req.Offset).Limit(req.Offset).Scan(&model)

	for key, item := range model {
		userName := fmt.Sprintf("%s:%d", "user", item.UserId)
		model[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(userName)
	}

	return gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": model}
}

func (ServiceRoom) ListByServiceManager(ServiceManagerId int, req Request.ServiceRoomList) gin.H {
	var model []Response.ServiceRoomList
	timeWhere := " where 1= 1"
	if req.StartTime != "" {
		timeWhere = fmt.Sprintf(" and users.create_time >= '%s'", req.StartTime)
	}

	if req.EndTime != "" {
		timeWhere += fmt.Sprintf(" and users.create_time <= '%s'", req.EndTime)
	}

	if req.UserName != "" {
		timeWhere += fmt.Sprintf(" and users.user_name = '%s'", req.UserName)
	}

	if req.ServiceName != "" {
		timeWhere += fmt.Sprintf(" and services.name = '%s'", req.ServiceName)
	}

	if req.ServiceMember != "" {
		timeWhere += fmt.Sprintf(" and services.username = '%s'", req.ServiceMember)
	}

	sql := `
			SELECT  user_name,user_head,
                         service_room_details.*,services.name,services.head as service_head
			FROM  users
				  left join service_room_details on users.user_id =  service_room_details.user_id
				  left join services on service_room_details.service_id = services.service_id
		`
	if ServiceManagerId != 0 {
		sql = fmt.Sprintf(`
			SELECT  user_name,user_head,
                         service_room_details.*,services.name,services.head as service_head
			FROM  (
					  select * from users where user_id in (
						  select * from (
							   select user_id from service_rooms where service_id in (
								   select service_id from services where service_manager_id = %d
							   )
						  ) a
					  )
				  ) as users
				  left join service_room_details on users.user_id =  service_room_details.user_id
				  left join services on service_room_details.service_id = services.service_id
		`, ServiceManagerId)
	}

	sql = sql + timeWhere
	tel := Base.MysqlConn.Raw(sql)

	// 计算分页和总数
	var allCount int
	telCount := Base.MysqlConn.Table("users")
	if ServiceManagerId != 0 {
		telCount = telCount.Where(`
		user_id in (
			select user_id from service_rooms where service_id in (
				select service_id from services where service_manager_id = ?
			)
		)
	`, ServiceManagerId)
	}
	type CountPage struct {
		Cnt int
	}
	var CountPageCnt CountPage
	sql = fmt.Sprintf("select count(*) as cnt from (%s) tt", sql)
	Base.MysqlConn.Raw(sql).Scan(&CountPageCnt)
	allCount = CountPageCnt.Cnt
	allPage := math.Ceil(float64(allCount) / float64(req.Offset))

	// 获取分页数据
	tel.Offset((req.Page - 1) * req.Offset).Limit(req.Offset).Scan(&model)

	for key, item := range model {
		model[key].IsOnline = Base.WebsocketHub.UserIdIsOnline(Common.Tools{}.GetUserWebSocketId(item.UserId))
	}

	return gin.H{"count": allCount, "page": allPage, "current_page": req.Page, "list": model}
}
