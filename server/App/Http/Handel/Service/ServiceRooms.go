package Service

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/Count"
	Message2 "server/App/Model/Message"
	Service2 "server/App/Model/Service"
	"server/App/Model/User"
	"server/Base"
	"time"
)

type ServiceRooms struct{}

// @summary 房间-获取用户房间列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_name query string false "用户名"
// @Param is_clear_status query string false "1清理count状态  0不清理"
// @Param type query string false "all 所有 user_no_read  用户未读 server_read 已回复 server_no_read 未回复 top 置顶 black 拉黑"
// @Param page query int true "指定页"
// @Param offset query int true "指定每页数量"
// @Router /service/rooms/list [post]
func (ServiceRooms) List(c *gin.Context) {
	var req Request.ServiceRoomList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetRoleId(c)
	res := Logic.ServiceRoom{}.List(RoleId, req, true)
	if req.IsClearStatus == 1 {
		Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), 0)
	}

	//序列化
	Common.ApiResponse{}.Success(c, "ok", res)
}

// @summary 房间-更新用户登录信息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "指定ID"
// @Param mobile query int true "手机号"
// @Param tag query int true "备注"
// @Param rename query int true "备注名字"
// @Router /service/rooms/end [post]
func (ServiceRooms) End(c *gin.Context) {
	var req Request.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), 0)
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", Common.Tools{}.GetRoleId(c), req.UserId).Updates(gin.H{"is_delete": 1, "late_msg": ""})
	Base.MysqlConn.Delete(
		&Message2.Message{}, "room_id =  ? and user_id = ?", Common.Tools{}.ConvertUserMessageRoomId(Common.Tools{}.GetRoleId(c), req.UserId), req.UserId)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"req": req})
}

// @summary 房间-更新用户登录信息
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "指定ID"
// @Param mobile query int true "手机号"
// @Param tag query int true "备注"
// @Param rename query int true "备注名字"
// @Router /service/rooms/update [post]
func (ServiceRooms) Update(c *gin.Context) {
	var req Request.ServiceRoomDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{"req": req})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)

	Logic.Service{}.Rename(serviceId, req.UserId, req.Rename)

	var findInfo Service2.ServiceRoomDetail
	Base.MysqlConn.Model(&findInfo).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Find(&findInfo)

	updates := gin.H{
		"mobile": req.Mobile,
		"tag":    req.Tag,
	}
	var info Service2.ServiceRoomDetail
	Base.MysqlConn.Model(&info).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Updates(updates)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"info": info})
}

// @summary 房间-获取用户房间详细
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Router /service/rooms/detail [post]
func (ServiceRooms) Detail(c *gin.Context) {
	var req Request.UserId
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	serviceId := Common.Tools{}.GetServiceId(c)
	Base.WebsocketHub.BindUser(Common.Tools{}.GetWebSocketId(c), req.UserId)
	Base.WebsocketHub.BindGroup(Common.Tools{}.GetWebSocketId(c), 0)

	// service未读更新0
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and user_id = ?", serviceId, req.UserId).Updates(gin.H{"service_no_read": 0})

	// 获取用户信息
	var users Response.UserDetail
	sql := "select service_rooms.`rename`,users.user_id,users.user_name,users.user_head,is_top,service_room_details.drive,service_room_details.ip,service_room_details.map,service_room_details.mobile,service_room_details.tag " +
		"from service_rooms left join users on service_rooms.user_id = users.user_id " +
		"left join service_room_details on service_rooms.room_id = service_room_details.room_id " +
		"where service_rooms.service_id = ? and service_rooms.user_id = ?"
	Base.MysqlConn.Raw(sql, serviceId, req.UserId).Scan(&users)

	var UserLoginRecorder []User.UserLoginLog
	Base.MysqlConn.Find(&UserLoginRecorder, "service_id = ? and user_id = ?", serviceId, req.UserId)

	var UserLoginRecorderResp []Response.UserLoginLog
	for _, v := range UserLoginRecorder {
		UserLoginRecorderResp = append(UserLoginRecorderResp, Response.UserLoginLog{
			Id:         v.Id,
			UserId:     v.UserId,
			ServiceId:  v.ServiceId,
			Ip:         v.Ip,
			Addr:       v.Addr,
			CreateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{"user": users, "login_log": UserLoginRecorderResp})
}

// @summary 房间-置顶指定用户
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Param top query int false "1置顶 0不置顶"
// @Router /service/rooms/top [post]
func (ServiceRooms) Top(c *gin.Context) {
	var req Request.RoomTop
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	serviceId := Common.Tools{}.GetServiceId(c)
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).
		Where("service_id = ? and user_id = ?", serviceId, req.UserId).Update("is_top", req.Top)
	var str string
	if req.Top == 1 {
		str = "置顶成功"
	} else {
		str = "取消置顶"
	}
	Common.ApiResponse{}.Success(c, str, gin.H{})
}

// @summary 房间-拉黑指定用户
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query int true "用户ID"
// @Param type query string true "拉黑类型ip user"
// @Param is_black query int false "1拉黑 0不拉黑"
// @Router /service/rooms/black [post]
func (ServiceRooms) Black(c *gin.Context) {
	var req Request.RoomBlack
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	serviceId := Common.Tools{}.GetServiceId(c)
	serviceModel, _ := Logic.Service{}.IdGet(serviceId)
	RoomId := Common.Tools{}.ConvertUserMessageRoomId(serviceId, req.UserId)
	if req.Type == "ip" {
		Logic.ServiceBlack{}.IpBlack(req.IsBlack, serviceId, req.Ip, serviceModel.ServiceManagerId)
	} else {
		Logic.ServiceBlack{}.UserBlack(req.IsBlack, serviceId, req.Ip, req.UserId, serviceModel.ServiceManagerId)
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{"room_id": RoomId})
}

// @summary 房间-拉黑列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param page query int true "指定页"
// @Param offset query int true "指定每页数量"
// @Router /service/rooms/black_list [post]
func (ServiceRooms) BlackList(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	tel := Base.MysqlConn.Table("service_blacks").
		Select("service_blacks.*,users.user_name,users.user_head").
		Joins("left JOIN users ON service_blacks.user_id = users.user_id").
		Where("service_blacks.service_id = ?", ServiceId)

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	var list []Response.UserBlackList
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}

func (ServiceRooms) Count(c *gin.Context) {
	serviceId := Common.Tools{}.GetServiceId(c)
	var list []Count.CountServiceRoom
	Base.MysqlConn.Order("id desc").Limit(5).Find(&list, "service_id = ?", serviceId)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"list": list})
}

// @summary 房间-备注
// @tags 客服系统
// @Param token header string true "认证token"
// @Param user_id query string true "用户ID"
// @Param rename query string false "备注"
// @Router /service/rooms/rename [post]
func (ServiceRooms) Rename(c *gin.Context) {
	var req Request.UpdateServiceRoomRename
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	Logic.Service{}.Rename(ServiceId, req.UserId, req.Rename)
	Common.ApiResponse{}.Success(c, "修改备注成功", gin.H{})
}

// @summary 房间-删除指定天前的用户
// @tags 客服系统
// @Param token header string true "认证token"
// @Param day query int false "天数前"
// @Router /service/rooms/delete_day [post]
func (ServiceRooms) DeleteDay(c *gin.Context) {
	var req Response.DeleteUserDay
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	ServiceId := Common.Tools{}.GetServiceId(c)

	now := time.Now()
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ? and update_time >= ?", ServiceId, now.AddDate(0, 0, -req.Day)).Updates(gin.H{"is_delete": 1, "late_msg": ""})
	Base.MysqlConn.Delete(&Message2.Message{}, "service_id = ? and create_time >= ?", ServiceId, now.AddDate(0, 0, -req.Day))

	Common.ApiResponse{}.Success(c, "用户清理成功", gin.H{"req": req, "time": now.AddDate(0, 0, -req.Day)})
}
