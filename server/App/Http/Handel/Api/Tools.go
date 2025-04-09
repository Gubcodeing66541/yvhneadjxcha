package Api

import (
	"fmt"
	"server/App/Common"
	"server/App/Http/Logic"
	Common2 "server/App/Model/Common"
	"server/App/Model/Service"
	"server/Base"
	"time"

	"github.com/gin-gonic/gin"
)

type Tools struct{}

// 话术复制
func (Tools) Copy(c *gin.Context) {
	var req struct {
		MyUsername   string   `json:"my_username"`
		CpoyUsername []string `json:"cpoy_username"`
		Selectd      []string `json:"selectd"` // head name hello quirk leave group
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	myService := Service.Service{}
	Base.MysqlConn.Find(&myService, "username = ?", req.MyUsername)
	if myService.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "我的账号不存在", gin.H{})
		return
	}

	//  拷贝账号
	serverList := []Service.Service{}
	for _, username := range req.CpoyUsername {
		copyService := Service.Service{}
		Base.MysqlConn.Find(&copyService, "username = ?", username)
		if copyService.ServiceId == 0 {
			Common.ApiResponse{}.Error(c, "复制账号不存在"+username, gin.H{})
			return
		}
		serverList = append(serverList, copyService)
	}

	//  拷贝密钥
	sm := []Service.ServiceMessage{}
	Base.MysqlConn.Find(&sm, "service_id = ? and `type` in ?", myService.ServiceId, req.Selectd)

	// 执行复制操作
	for _, server := range serverList {
		if isInArray(req.Selectd, "head") {
			server.Head = myService.Head
		}
		if isInArray(req.Selectd, "name") {
			server.Name = myService.Name
		}

		Base.MysqlConn.Save(&server)

		for _, h := range sm {
			Base.MysqlConn.Create(&Service.ServiceMessage{
				ServiceId:  server.ServiceId,
				Type:       h.Type,
				MsgType:    h.MsgType,
				MsgInfo:    h.MsgInfo,
				Status:     h.Status,
				CreateTime: time.Now(),
			})
		}
	}

	Common.ApiResponse{}.Success(c, "复制成功", gin.H{"sm": sm})
}

func isInArray(arr []string, target string) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

// code search
func (Tools) Search(c *gin.Context) {
	var req struct {
		Username []string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	type response struct {
		Username string `json:"username"`
		Status   string `json:"status"`
	}

	var serviceInfo Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&serviceInfo)
	if serviceInfo.Id == 0 {
		Common.ApiResponse{}.Error(c, "无法查询到卡密", gin.H{})
		return
	}

	domain := Logic.Domain{}.GetPublic()
	if serviceInfo.BindDomainId != 0 {
		domainTemp := Logic.Domain{}.Get(serviceInfo.BindDomainId)
		if domainTemp.Domain != "" && domainTemp.Status == "enable" {
			domain = domainTemp.Domain
		}
	}
	web := domain + "?code=" + serviceInfo.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix())
	if serviceInfo.BindDomain != "" {
		web = serviceInfo.BindDomain + "?code=" + serviceInfo.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix())
	}

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"domain": web})
}

func (Tools) ServiceCount(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&service)
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

func (Tools) GetCode(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&service)
	if service.Id == 0 {
		Common.ApiResponse{}.Error(c, "服务不存在", gin.H{})
		return
	}

	Logic.Domain{}.Bind(service.Id)

	var domainList []Common2.Domain
	Base.MysqlConn.Find(&domainList, "id in (?)", []int{service.BindDomainId, service.BindDomainId2, service.BindDomainId3})

	for i, domain := range domainList {
		web := domain.Domain + "?code=" + service.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix())
		domainList[i].Domain = web
	}

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"domain": domainList, "service": service})
}

func (Tools) FixDomain(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&service)
	if service.Id == 0 {
		Common.ApiResponse{}.Error(c, "服务不存在", gin.H{})
		return
	}

	// 绑定域名
	Logic.Domain{}.Bind(service.Id)

	var domainList []Common2.Domain
	Base.MysqlConn.Find(&domainList, "id in (?)", []int{service.BindDomainId, service.BindDomainId2, service.BindDomainId3})

	for i, domain := range domainList {
		web := domain.Domain + "?code=" + service.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix())
		domainList[i].Domain = web
	}

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"domain": domainList, "service": service})
}

func (Tools) ResetDomain(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&service)
	if service.Id == 0 {
		Common.ApiResponse{}.Error(c, "服务不存在", gin.H{})
		return
	}

	// 重设置code
	service.Code = Common.Tools{}.CreateActiveCode(service.ServiceId)
	Base.MysqlConn.Save(&service)

	var domainList []Common2.Domain
	Base.MysqlConn.Find(&domainList, "id in (?)", []int{service.BindDomainId, service.BindDomainId2, service.BindDomainId3})

	for i, domain := range domainList {
		web := domain.Domain + "?code=" + service.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix())
		domainList[i].Domain = web
	}

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"domain": domainList, "service": service})
}

func (Tools) Count(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
		return
	}

	var service Service.Service
	Base.MysqlConn.Where("username = ?", req.Username).Find(&service)
	if service.Id == 0 {
		Common.ApiResponse{}.Error(c, "服务不存在", gin.H{})
		return
	}

	var sql1Result []struct {
		UserCnt int64  `json:"user_cnt"`
		Dt      string `json:"dt"`
	}
	sql1 := "SELECT count(*) as user_cnt,DATE_FORMAT(create_time,'%Y/%m/%d') as dt  FROM `service_rooms` where service_id = ? GROUP by dt"
	Base.MysqlConn.Raw(sql1, service.Id).Scan(&sql1Result)

	var sql2Result []struct {
		IpCnt int64  `json:"ip_cnt"`
		Dt    string `json:"dt"`
	}
	sql2 := "SELECT count(DISTINCT late_id) as ip_cnt,DATE_FORMAT(create_time,'%Y/%m/%d') as dt FROM `service_rooms` where service_id = ? GROUP by dt"
	Base.MysqlConn.Raw(sql2, service.Id).Scan(&sql2Result)

	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"sql1": sql1Result, "sql2": sql2Result})
}
