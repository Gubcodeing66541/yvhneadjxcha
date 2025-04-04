package Api

import (
	"server/App/Common"
	"server/App/Model/Service"
	"server/Base"
	"time"

	"github.com/gin-gonic/gin"
)

type CommonHanel struct{}

func (CommonHanel) ServiceCount(c *gin.Context) {
	username := c.Param("username")

	var service Service.Service
	Base.MysqlConn.Where("user_name = ?", username).Find(&service)
	if service.Id == 0 {
		Common.ApiResponse{}.Error(c, "服务不存在", nil)
		return
	}

	// 创建最近7天的日期数组
	dates := make([]string, 7)
	ipCounts := make([]int64, 7)
	for i := 0; i < 7; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		dates[6-i] = date
	}

	// 查询最近7天的IP统计数据
	type Result struct {
		Date  string
		Count int64
	}
	var results []Result
	Base.MysqlConn.Model(&Service.ServiceRoom{}).
		Select("DATE_FORMAT(create_time, '%Y-%m-%d') as date, COUNT(DISTINCT late_ip) as count").
		Where("service_id = ?", service.Id).
		Where("create_time > ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE_FORMAT(create_time, '%Y-%m-%d')").
		Find(&results)

	// 将查询结果映射到对应日期
	for _, result := range results {
		for i, date := range dates {
			if date == result.Date {
				ipCounts[i] = result.Count
				break
			}
		}
	}

	// 构建返回数据
	responseData := make(map[string]interface{})
	responseData["dates"] = dates
	responseData["ip_counts"] = ipCounts

	Common.ApiResponse{}.Success(c, "获取成功", responseData)
}
