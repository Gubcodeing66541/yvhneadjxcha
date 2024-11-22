package Logic

import (
	"server/App/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Auth struct{}

func (Auth) Register(username string, password string, serviceManagerId int) int {
	times := time.Now()
	serviceAuth := &Service2.ServiceAuth{ServiceManagerId: serviceManagerId, Username: username, Password: password, CreateTime: times, UpdateTime: times, TimeOut: times}
	Base.MysqlConn.Create(&serviceAuth)
	code := Common.Tools{}.CreateActiveCode(serviceAuth.ServiceId)
	head := Base.AppConfig.HttpHost + "/static/static/head.png"
	Base.MysqlConn.Create(&Service2.Service{
		ServiceManagerId: serviceManagerId, IsActivate: 0, ActivateTime: time.Now(),
		ServiceId: serviceAuth.ServiceId, Name: "小客服", Head: head, Username: username,
		Code: code, CreateTime: times, Type: "auth", Role: "user", TimeOut: times, Status: "success",
	})

	Base.MysqlConn.Create(&Service2.ServiceNoticeSetting{ServiceId: serviceAuth.ServiceId})
	return serviceAuth.ServiceId
}

func (Auth) RegisterByServiceManager(username string, serviceName string, serviceManagerId int, Day int) int {
	times := time.Now()
	serviceAuth := &Service2.ServiceAuth{ServiceManagerId: serviceManagerId, Username: username, CreateTime: times, UpdateTime: times, TimeOut: times}
	Base.MysqlConn.Create(&serviceAuth)
	code := Common.Tools{}.CreateActiveCode(serviceAuth.ServiceId)
	head := Base.AppConfig.HttpHost + "/static/static/head.png"
	Base.MysqlConn.Create(&Service2.Service{
		ServiceManagerId: serviceManagerId, IsActivate: 0, Day: 0, ActivateTime: time.Now(), Status: "success",
		ServiceId: serviceAuth.ServiceId, Name: serviceName, Head: head, Username: username,
		Code: code, CreateTime: times, Type: "auth", Role: "user", TimeOut: times, CodeBackground: "#ffffff", CodeColor: "#000000",
	})
	return serviceAuth.ServiceId
}
