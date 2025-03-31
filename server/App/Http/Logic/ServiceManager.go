package Logic

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"server/App/Common"
	"server/App/Http/Response"
	ServiceManager2 "server/App/Model/ServiceManager"
	"server/Base"
	"time"
)

type ServiceManager struct{}

func (ServiceManager) Get(serviceManagerId int) ServiceManager2.ServiceManager {
	var model ServiceManager2.ServiceManager
	Base.MysqlConn.Find(&model, "service_manager_id = ?", serviceManagerId)
	return model
}

func (ServiceManager) Create(c *gin.Context, Account int) (username string, password string) {
	username = Common.Tools{}.CreateServiceManagerMember()
	password = "abc123456"
	manager := ServiceManager2.ServiceManager{
		Member:     username,
		Password:   password,
		Name:       "小客服",
		Head:       Common.Tools{}.GetDefaultHead(),
		Ip:         "-",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Account:    Account,
		Status:     "success",
	}
	Base.MysqlConn.Create(&manager)

	// 生成机器人
	Base.MysqlConn.Create(&ServiceManager2.ServiceManagerBot{
		ServiceManagerId: manager.ServiceManagerId,
		Status:           "stop",
		Head:             Base.AppConfig.HttpHost + "/static/static/bot.png",
		Name:             "机器人",
	})
	return username, password
}

func (ServiceManager) Renew(ServiceManagerId int, Account int, Reason string, PayType string, member string) {
	var account ServiceManager2.ServiceManager
	Base.MysqlConn.Find(&account, "service_manager_id = ?", ServiceManagerId)
	if account.ServiceManagerId == 0 {
		return
	}

	if member == "" {
		member = account.Member
	}

	Base.MysqlConn.Model(&ServiceManager2.ServiceManager{}).Where("service_manager_id = ?", ServiceManagerId).
		Update("account", gorm.Expr("account + ?", Account))
	Base.MysqlConn.Create(&Response.ServiceManagerRenewRecorder{
		ServiceManagerMember: account.Member, Member: member,
		ServiceManagerId: ServiceManagerId, OldAccount: account.Account, OrderId: Common.Tools{}.CreateOrderId("PAY"),
		Account: account.Account + Account, Renew: Account, CreateTime: time.Now().Format("2006-01-02 15:04:05"), Reason: Reason, PayType: PayType,
	})
}

func (ServiceManager) RenewByMember(Account int, Reason string, PayType string, member string) {
	var account ServiceManager2.ServiceManager
	Base.MysqlConn.Find(&account, "member = ?", member)
	if account.ServiceManagerId == 0 {
		return
	}

	Base.MysqlConn.Model(&ServiceManager2.ServiceManager{}).Where("service_manager_id = ?", account.ServiceManagerId).
		Update("account", gorm.Expr("account + ?", Account))
	Base.MysqlConn.Create(&Response.ServiceManagerRenewRecorder{
		ServiceManagerMember: account.Member, Member: member,
		ServiceManagerId: account.ServiceManagerId, OldAccount: account.Account, OrderId: Common.Tools{}.CreateOrderId("PAY"),
		Account: account.Account + Account, Renew: Account, CreateTime: time.Now().Format("2006-01-02 15:04:05"), Reason: Reason, PayType: PayType,
	})
}
