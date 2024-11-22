package Servicemanager

import (
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common2 "server/App/Model/Common"
	"server/App/Model/Count"
	"server/App/Model/Service"
	"server/App/Model/ServiceManager"
	"server/Base"
	"strings"
)

type Manager struct {
}

// @summary 基础-获取账号信息详情
// @tags 客服后台
// @Description service_manager_id 账号ID
// @Description account 钱包
// @Description id_card_a身份证A面
// @Description id_card_b身份证B面
// @Description head头像
// @Description real_name 真名
// @Description account 钱包
// @Description head头像
// @Description real_name 真名
// @Description id_card_number 身份证号
// @Description id_card_number 身份证号
// @Description id_card_a 身份证A面
// @Description id_card_b 身份证B面
// @Description id_card_apply_status 身份证状态 none没有上传  apply 审批中 agree 同意 refuse 拒绝
// @Param token header string true "认证token"
// @Router /service/manager/info [post]
func (Manager) Info(c *gin.Context) {
	Id := Common.Tools{}.GetRoleId(c)
	var model Response.ServiceManagerInfo
	Base.MysqlConn.Find(&ServiceManager.ServiceManager{}, "service_manager_id = ?", Id).Scan(&model)

	if model.ServiceManagerId == 0 || model.Status == "no_use" {
		Common.ApiResponse{}.NoAuth(c, gin.H{})
		return
	}

	var config Response.ResConfig
	Base.MysqlConn.Model(&Common2.SystemConfig{}).Scan(&config)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"info": model, "role_id": Id, "config": config})
}

// @summary 基础-获取统计信息
// @tags 客服后台
// @Param token header string true "认证token"
// @Router /service/manager/count [post]
func (Manager) Count(c *gin.Context) {
	Id := Common.Tools{}.GetRoleId(c)
	var model Response.ServiceManagerCount
	Base.MysqlConn.Raw("select "+
		"(select count(*) from service_rooms where create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00') and service_id in ((select service_id from services where service_manager_id = ?))) as room_cnt,"+
		"(select count(*) from messages  where create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00') and service_id in (select service_id from services where service_manager_id = ?)) as message_cnt",
		Id, Id).Scan(&model)

	// 查询并排序
	var CountRecorderTemp []Count.CountServiceRoom
	var CountRecorder []Count.CountServiceRoom
	Base.MysqlConn.Limit(30).Order("id desc").Find(&CountRecorderTemp, "service_manager_id = ? and service_id = 0", Id)

	for i := len(CountRecorderTemp); i > 0; i-- {
		CountRecorder = append(CountRecorder, CountRecorderTemp[i-1])
	}
	userOnlineCnt := make([]int, 0)
	serviceOnlineCnt := make([]int, 0)
	timeKey := make([]string, 0)

	for key, _ := range CountRecorder {
		userOnlineCnt = append(userOnlineCnt, CountRecorder[key].OnlineUser)
		serviceOnlineCnt = append(serviceOnlineCnt, CountRecorder[key].OnlineService)
		timeKey = append(timeKey, CountRecorder[key].CountTime.Format("2006-01-02"))
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{
		"count":          model,
		"service_online": serviceOnlineCnt,
		"user_online":    userOnlineCnt,
		"time_key":       timeKey,
	})
}

// @summary 基础-修改账号基本信息接口
// @tags 客服后台
// @Param token header string true "认证token"
// @Param name query string false "昵称"
// @Param head query string false "头像"
// @Param id_card_a query string false "身份证a"
// @Param id_card_b query string false "身份证b"
// @Param real_name query string false "真实名称"
// @Param id_card_number query string false "身份证号码"
// @Router /service/manager/update [post]
func (Manager) Update(c *gin.Context) {
	var req Request.ServiceManagerUpdate
	if req.IdCardNumber != "" {
		req.IdCardApplyStatus = "apply"
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数填写有误，请重新填写", gin.H{})
		return
	}
	Id := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&ServiceManager.ServiceManager{}).Where("service_manager_id = ?", Id).Updates(req)
	Common.ApiResponse{}.Success(c, "保存成功", gin.H{})
}

// @summary 基础-获取统计信息用户
// @tags 客服后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/count_room_detail [post]
func (Manager) CountRoomDetail(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}
	Id := Common.Tools{}.GetRoleId(c)
	var model []Response.ServiceManagerRoomCount

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Raw(""+
		"select service_id,count(*) as row_count from service_rooms where service_id in ("+
		"select service_id from services where service_manager_id = ? and  create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00')"+
		") GROUP BY service_id", Id).Scan(&model)

	// 计算分页和总数
	var allCountInterface interface{}
	Base.MysqlConn.Raw(""+
		"select count(service_id) as row_count from service_rooms where service_id in ("+
		"select service_id from services where service_manager_id = ? and  create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00') "+
		") GROUP BY service_id", Id).Count(&allCountInterface)

	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	Base.MysqlConn.Raw("select t.*,services.name from ("+
		"select service_id,count(*) as room_cnt from service_rooms where service_id in ("+
		"select service_id from services where service_manager_id = ?"+
		") and  create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00')  GROUP BY service_id limit ? offset ?) t left join services on t.service_id = services.service_id", Id, pageReq.Offset, (pageReq.Page-1)*pageReq.Offset).Scan(&model)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"list": model, "count": len(model), "page": allPage, "current_page": pageReq.Page})
}

// @summary 基础-获取统计信息消息
// @tags 客服后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/manager/count_message_detail [post]
func (Manager) CountMessageDetail(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}
	Id := Common.Tools{}.GetRoleId(c)
	// 计算分页和总数
	var allCount []Response.Count
	Base.MysqlConn.Raw("select service_id,count(*) as counted  from messages where service_id in ("+
		"select service_id from services where service_manager_id = ?"+
		") GROUP BY service_id", Id).Scan(&allCount)

	allPage := math.Ceil(float64(len(allCount)) / float64(pageReq.Offset))

	var model []Response.ServiceManagerMessageCount
	Base.MysqlConn.Raw("select t.*,services.name from (select service_id,count(*) as message_cnt from messages where service_id in ("+
		"select service_id from services where service_manager_id = ?"+
		") and  create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00')  GROUP BY service_id limit ? offset ?) t left join services on t.service_id = services.service_id", Id, pageReq.Offset, (pageReq.Page-1)*pageReq.Offset).Scan(&model)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"list": model, "count": len(allCount), "page": allPage, "current_page": pageReq.Page})
}

// @summary 基础-获取数据大屏
// @tags 客服后台
// @Param token header string true "认证token"
// @Router /service/manager/data [post]
func (Manager) Data(c *gin.Context) {
	roleId := Common.Tools{}.GetRoleId(c)
	serviceIds, UserIds := Common.Socket{}.GetAll()
	var socketDataResponse Response.SocketDataResponse
	sql := "select " +
		"(select count(*) from service_rooms where service_id in (select service_id from services where service_manager_id = ?)) as all_user_cnt," +
		"(select count(*) from messages where service_id in (select service_id from services where service_manager_id = ?)) as all_message_cnt," +
		" (select count(*) from service_rooms where create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00') and service_id in (select  service_id from services where service_manager_id = ?))  as today_user_cnt," +
		"(select count(*) from messages where create_time >= DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00') and service_id in (select  service_id from services where service_manager_id = ?)) as today_message_cnt," +
		"(select count(*) from service_rooms  where service_id in (select  service_id from services where service_manager_id = ?) and user_id in (?)) as online_user," +
		"(select count(*) from services where service_id in (select  service_id from services where service_manager_id = ? and service_id in (?))) as online_service"
	Base.MysqlConn.Raw(sql, roleId, roleId, roleId, roleId, roleId, UserIds, roleId, serviceIds).Scan(&socketDataResponse)

	// 大屏中间的区域统计
	var regionResponse []Response.RegionResponse
	sql = "select map,count(*) as cnt from service_room_details where service_id in (select service_id from services where service_manager_id = ?) group by map"
	Base.MysqlConn.Raw(sql, roleId).Scan(&regionResponse)

	// 在线用户统计
	var OnlineUser []Response.SocketUserResponse
	sql = "select user_name,user_head,map,services.name as service_name from service_rooms " +
		"left join users on service_rooms.user_id = users.user_id " +
		"left join services on service_rooms.service_id = services.service_id " +
		"left join service_room_details on service_rooms.service_id = service_room_details.service_id and service_rooms.user_id = service_room_details.user_id " +
		"where service_rooms.service_id in (select service_id from services where service_manager_id = ?) and service_rooms.user_id in (?)"
	Base.MysqlConn.Raw(sql, roleId, UserIds).Scan(&OnlineUser)

	// 在线客服统计
	var serviceOnline []Service.Service
	sql = "select * from services where service_manager_id = ? and service_id in (?)"
	Base.MysqlConn.Raw(sql, roleId, serviceIds).Scan(&serviceOnline)

	// 客服访问排名
	var serviceRank []Response.ServiceRank
	sql = "select t1.*,services.name,services.head from (" +
		"select service_id,count(*) as cnt from service_rooms where service_id in (" +
		"select  service_id from services where service_manager_id = ?" +
		") group by service_id" +
		") t1 left join services on t1.service_id = services.service_id order by cnt desc"
	Base.MysqlConn.Raw(sql, roleId).Scan(&serviceRank)

	regionMap := make(map[string]int)
	for _, item := range regionResponse {
		regionName := item.Map
		if strings.Index(regionName, "省") >= 0 {
			regionName = strings.Split(regionName, "省")[0]
		} else if strings.Index(regionName, "市") >= 0 {
			regionName = strings.Split(regionName, "市")[0]
		} else {
			continue
		}

		if _, ok := regionMap[regionName]; !ok {
			regionMap[regionName] = 0
		}
		regionMap[regionName] = regionMap[regionName] + item.Cnt
	}

	regionResponse = make([]Response.RegionResponse, 0)
	for regionKey, regionCount := range regionMap {
		regionResponse = append(regionResponse, Response.RegionResponse{
			Map: regionKey,
			Cnt: regionCount,
		})
	}

	var contTime []Response.ContTime
	sql = "select DATE_FORMAT(create_time,'%Y-%m-%d') as time ,count(*) as cnt from service_rooms " +
		"where service_id in (select service_id from services where service_manager_id = ?) " +
		"group by DATE_FORMAT(create_time,'%Y-%m-%d') order by time desc limit 7"
	Base.MysqlConn.Raw(sql, roleId).Scan(&contTime)

	Common.ApiResponse{}.Success(c, "ok", gin.H{
		"temp":           gin.H{"user": UserIds, "service": serviceIds},
		"top":            socketDataResponse,
		"region":         regionResponse,
		"user_online":    OnlineUser,
		"service_online": serviceOnline,
		"service_rank":   serviceRank,
		"count_time":     contTime,
	})
}
