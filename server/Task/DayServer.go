package Task

import (
	"encoding/json"
	"fmt"
	"server/App/Common"
	Count2 "server/App/Model/Count"
	Service2 "server/App/Model/Service"
	ServiceManager2 "server/App/Model/ServiceManager"
	"server/Base"
	"time"
)

type Day struct{}

type Count struct {
	AddUser       int `json:"add_user"`
	AllUser       int `json:"all_user"`
	OnlineUser    int `json:"online_user"`
	OnlineService int `json:"online_service"`
}

type Req struct {
	Data ReqData `json:"data"`
}
type ReqData struct {
	ServiceIds []int `json:"service_ids"`
	UserIds    []int `json:"user_ids"`
}

func (d Day) Run() {
	fmt.Println("")
	fmt.Println("执行本次统计任务", time.Now())

	var req Req
	reqStr := Common.Tools{}.HttpGet("http://127.0.0.1/api/system/status_all")
	err := json.Unmarshal(reqStr, &req)
	if err != nil {
		fmt.Println("在线统计请求失败", reqStr)
		return
	}

	fmt.Println("在线数据", req)
	fmt.Println("serviceIds", req.Data.ServiceIds)
	fmt.Println("userIds", req.Data.UserIds)

	// 统计所有的代理
	var serviceManagerList []ServiceManager2.ServiceManager
	Base.MysqlConn.Find(&serviceManagerList)

	// 对单个代理执行统计
	for _, item := range serviceManagerList {
		d.CountServiceManager(item.ServiceManagerId, req.Data.ServiceIds, req.Data.UserIds)
	}

}

// CountServiceManager 统计代理账号
func (d Day) CountServiceManager(ServiceManagerId int, serviceIds []int, userIds []int) {
	fmt.Println("统计代理:", ServiceManagerId)

	var service []Service2.Service
	Base.MysqlConn.Find(&service, "service_manager_id = ?", ServiceManagerId)

	now := time.Now()
	var count Count
	sql := "select" +
		"(select count(*) from service_rooms where service_id in (select service_id from services where service_manager_id = ? ) ) as add_user," +
		"(select count(*) from service_rooms  where service_id in (select service_id from services where service_manager_id = ?) ) as all_user," +
		"(select count(*) from service_rooms  where service_id in (select service_id from services where service_manager_id = ?) and user_id in (?)) as online_user," +
		"(select count(*) from services where service_id in (select service_id from services where service_manager_id = ? and service_id in (?))) as online_service"

	Base.MysqlConn.Raw(sql, ServiceManagerId, ServiceManagerId, ServiceManagerId, userIds, ServiceManagerId, serviceIds).Scan(&count)

	Base.MysqlConn.Create(&Count2.CountServiceRoom{
		ServiceManagerId: ServiceManagerId,
		ServiceId:        0,
		AllUser:          count.AllUser,
		AddUser:          count.AddUser,
		OnlineUser:       count.OnlineUser,
		OnlineService:    count.OnlineService,
		CreateTime:       now,
		CountTime:        now,
	})

	for _, item := range service {
		d.CountService(ServiceManagerId, item.ServiceId, userIds)
	}
}

// CountService 统计客服账号
func (d Day) CountService(ServiceManagerId int, serviceId int, userIds []int) {
	fmt.Println("统计客服:", serviceId)

	now := time.Now()
	var count Count
	sql := "select" +
		"(select count(*) from service_rooms where service_id  = ? ) as add_user," +
		"(select count(*) from service_rooms  where service_id = ? ) as all_user," +
		"(select count(*) from service_rooms  where service_id = ? and user_id in (?)) as online_user"
	Base.MysqlConn.Raw(sql, serviceId, serviceId, serviceId, userIds).Scan(&count)

	Base.MysqlConn.Create(&Count2.CountServiceRoom{
		ServiceManagerId: ServiceManagerId,
		ServiceId:        serviceId,
		AllUser:          count.AllUser,
		AddUser:          count.AddUser,
		OnlineUser:       count.OnlineUser,
		OnlineService:    count.OnlineService,
		CreateTime:       now,
		CountTime:        now,
	})
}
