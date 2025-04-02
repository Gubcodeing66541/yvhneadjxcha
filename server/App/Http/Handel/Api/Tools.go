package Api

import (
	"server/App/Common"
	"server/App/Model/Service"
	"server/Base"
	"time"

	"github.com/gin-gonic/gin"
)

type Tools struct{}

func (Tools) Copy(c *gin.Context) {
	var req struct {
		MyUsername   string   `json:"my_username"`
		CpoyUsername []string `json:"cpoy_username"`
		Selectd      []string `json:"selectd"` // head name hello quirk leave group
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
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
	var sm []Service.ServiceMessage
	Base.MysqlConn.Find(&sm, "service_id = ? and `type` in ?", myService.ServiceId, req.Selectd)

	// 执行复制操作
	for _, server := range serverList {
		if isInArray(req.Selectd, "head") {
			server.Head = myService.Head
		}
		if isInArray(req.Selectd, "name") {
			server.Name = myService.Name
		}

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

		Base.MysqlConn.Save(&server)
	}

}

func isInArray(arr []string, target string) bool {
	for _, num := range arr {
		if num == target {
			return true
		}
	}
	return false
}

func (Tools) Get(c *gin.Context) {
	var req struct {
		Username []string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
	}

	type response struct {
		Username string `json:"username"`
		Timeout  int64  `json:"timeout"`
	}

}

func (Tools) Search(c *gin.Context) {
	var req struct {
		Username []string `json:"username"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数", gin.H{})
	}

	type response struct {
		Username string `json:"username"`
		Status   string `json:"status"`
	}

}
