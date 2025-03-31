package Logic

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"server/App/Common"
	"server/App/Http/Dto"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Service struct{}

// code获取客服信息
func (s Service) Get(code string) (Service2.Service, error) {
	var model Service2.Service
	Base.MysqlConn.Find(&model, "code = ? ", code)
	return model, nil
}

func (Service) List(req Request.GerServiceList) Response.ResultList {
	var result Response.ResultList
	page := req.Page
	result.Page = page

	var respServiceList []Response.RespServiceList
	if req.Name != "" {
		Base.MysqlConn.Raw("select service.service_id,service.name,service.head,service.username,service.type,service.role,service.code,service.time_out,"+
			"domain.domain,domain.status from services as service left join domains as domain on service.service_id = domain.bind_service_id  "+
			"where service.name like ? or service.username like ? order by service.id desc limit ?,?   ",
			"%"+req.Name+"%", "%"+req.Name+"%",
			(req.Page.CurrentPage-1)*req.Page.CurrentSize,
			req.Page.CurrentSize).Scan(&respServiceList)
	} else {
		Base.MysqlConn.Raw("select service.service_id,service.name,service.head,service.username,service.type,service.role,service.code,service.time_out,"+
			"domain.domain,domain.status from services as service left join domains as domain on service.service_id = domain.bind_service_id  "+
			"order by service.id desc limit ?,?   ",
			(req.Page.CurrentPage-1)*req.Page.CurrentSize,
			req.Page.CurrentSize).Scan(&respServiceList)
	}
	result.Data = respServiceList

	var count int
	Base.MysqlConn.Model(&Service2.Service{}).Count(&count)
	result.Page.Count = count
	return result
}

func (Service) GetServiceDomain(serviceId int) string {
	domain := Domain{}.GetServiceBind(serviceId)
	return domain.Domain
}

func (s Service) IdGet(serviceId int) (Service2.Service, error) {
	var model Service2.Service
	Base.MysqlConn.Find(&model, "service_id = ?", serviceId)
	return model, nil
}

// 系统注册账号
func (Service) Create() (serviceId int, member string, password string) {
	member = Common.Tools{}.CreateActiveMember()
	password = "Ab147258369"
	serviceId = Auth{}.Register(member, password, 0)
	return serviceId, member, password
}

func (Service) BachCreate(req Request.CreateServiceDay) []Dto.ServiceLoginInfo {
	var ServiceLoginInfos []Dto.ServiceLoginInfo
	for i := 0; i < req.Count; i++ {
		member := Common.Tools{}.CreateActiveMember()
		password := "Ab147258369"
		serviceId := Auth{}.Register(member, password, 0)
		ServiceLoginInfos = append(ServiceLoginInfos, Dto.ServiceLoginInfo{ServiceId: serviceId, Member: member, Password: password})
	}
	return ServiceLoginInfos
}

// 发送消息给用户
func (Service) SendMessageToUser(serviceId int, userId int, msgType string, message string) {
	roomId := Common.Tools{}.ConvertUserMessageRoomId(serviceId, userId)
	fmt.Println(roomId)
}

func (s Service) Renewal(serviceId int, day int) error {
	//查询客服是否存在
	var model Service2.ServiceAuth
	Base.MysqlConn.Find(&model, "service_id = ?", serviceId)
	if model.ServiceId == 0 {
		return errors.New("当前客服不存在")
	}

	//生成订单
	Order{}.Create(serviceId, day)

	//如果是过期用户直接加上天数
	now := time.Now()
	model.UpdateTime = now

	//更新过期时间
	var timeOut time.Time
	if model.TimeOut.Unix()-time.Now().Unix() <= 0 {
		timeOut = now.AddDate(0, 0, day)
	} else {
		timeOut = model.TimeOut.AddDate(0, 0, day)
	}
	model.TimeOut = timeOut
	Base.MysqlConn.Save(&model)

	//客服升级为VIP
	var service Service2.Service
	update := gin.H{"role": "vip", "time_out": timeOut}
	Base.MysqlConn.Model(&service).Where("service_id = ?", serviceId).Updates(update)

	// 绑定域名
	_ = Domain{}.Bind(model.ServiceId)

	// 清理缓存
	s.ClearCache(model.ServiceId)
	return nil
}

func (s Service) RenewalByUsername(username string, day int) error {
	//查询客服是否存在
	var model Service2.ServiceAuth
	Base.MysqlConn.Find(&model, "username = ?", username)
	if model.ServiceId == 0 {
		return errors.New("当前客服不存在")
	}

	//生成订单
	Order{}.Create(model.ServiceId, day)

	//如果是过期用户直接加上天数
	now := time.Now()
	model.UpdateTime = now

	//更新过期时间
	var timeOut time.Time
	if model.TimeOut.Unix()-time.Now().Unix() <= 0 {
		timeOut = now.AddDate(0, 0, day)
	} else {
		timeOut = model.TimeOut.AddDate(0, 0, day)
	}
	model.TimeOut = timeOut
	Base.MysqlConn.Save(&model)

	//客服升级
	var service Service2.Service
	update := gin.H{"role": "vip", "time_out": timeOut}
	Base.MysqlConn.Model(&service).Where("service_id = ?", model.ServiceId).Updates(update)

	if service.IsActivate == 0 {
		Base.MysqlConn.Model(&service).Where("service_id = ?", model.ServiceId).Update("day", gorm.Expr("day+ ?", day))
	}

	// 绑定域名
	_ = Domain{}.Bind(model.ServiceId)

	// 清理缓存
	s.ClearCache(model.ServiceId)
	return nil
}

func (s Service) RenewalByServiceManager(serviceManagerId int, username string, day int, reason string) error {

	// 检查余额
	serviceManager := ServiceManager{}.Get(serviceManagerId)
	pay := Order{}.GerAmount()
	if serviceManager.Account < pay*day {
		return errors.New("您的余额已不足")
	}

	// 执行扣费
	ServiceManager{}.Renew(serviceManagerId, -pay*day, reason, "system", username)

	// 续费账号
	return s.RenewalByUsername(username, day)
}

func (Service) IsMaturities(serviceId int) error {
	//查询客服是否存在
	var model Service2.ServiceAuth
	Base.MysqlConn.Find(&model, "service_id = ?", serviceId)
	if model.ServiceId == 0 {
		return errors.New("客服不存在")
	}

	//过期
	now := time.Now()
	if model.TimeOut.Before(now) {
		return errors.New("客服过期")
	}
	return nil
}

// 清除缓存
func (s Service) ClearCache(serviceId int) {
}

// 获取客服角色
func (s Service) GetServiceRole(serviceId int) string {
	service, err := s.IdGet(serviceId)
	if err != nil {
		return ""
	}
	return service.Role
}
func (Service) Rename(ServiceId int, userId int, rename string) {
	var serviceRoom Service2.ServiceRoom
	Base.MysqlConn.Model(&serviceRoom).Where("user_id = ? and service_id =?", userId, ServiceId).Update("rename", rename)
}

func (Service) DelService(req Request.DelService) {
	var serviceRoom Service2.ServiceRoom
	Base.MysqlConn.Model(&serviceRoom).Where("id = ? and service_id = ?", req.Id, req.ServiceId).Delete(serviceRoom)
}

// 系统注册账号
func (Service) CreateByServiceManager(serviceManagerId int, serviceName string, Day int) (serviceId int, member string, err error) {
	member = Common.Tools{}.CreateActiveMember()
	serviceId, err = Auth{}.RegisterByServiceManager(member, serviceName, serviceManagerId, Day)
	return serviceId, member, err
}
