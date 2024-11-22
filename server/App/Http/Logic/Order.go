package Logic

import (
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"time"
)

type Order struct {
}

// 每日续费单价
const DayRenewalAmount int = 70

func (Order) GerAmount() int {
	var setting Common.SystemConfig
	Base.MysqlConn.Find(&setting)
	if setting.Pay != 0 {
		return setting.Pay
	}
	return DayRenewalAmount
}

// ID查询订单
func (Order) QueryById(id int) Service2.Order {
	var entity Service2.Order
	Base.MysqlConn.Find(&entity, "id = ?", id)
	return entity
}

// 订单列表
func (Order) List() []Service2.Order {
	var entity []Service2.Order
	Base.MysqlConn.Find(&entity)
	if entity != nil {
		return entity
	}
	return nil
}

// 订单列表
func (Order) GetByServiceId(serviceId int) []Service2.Order {
	var entity []Service2.Order
	Base.MysqlConn.Find(&entity, "service_id = ?", serviceId)
	if entity != nil {
		return entity
	}
	return nil
}

// 创建订单
func (Order) Create(serviceId int, day int) Service2.Order {
	var order Service2.Order
	order.CreateTime = time.Now()
	order.ServiceId = serviceId
	order.Amount = day * Order{}.GerAmount()
	order.Day = day

	Base.MysqlConn.Create(&order)
	return order
}
