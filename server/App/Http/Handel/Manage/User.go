package Manage

import (
	"github.com/gin-gonic/gin"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Model/Message"
	"server/Base"
	"time"
)

type User struct{}

// @summary 用户列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/users/list [post]
func (User) List(c *gin.Context) {
	var req Request.ServiceRoomList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	res := Logic.ServiceRoom{}.ListByServiceManager(0, req)

	//序列化
	Common.ApiResponse{}.Success(c, "ok", res)
}

func (User) ClearMessage(c *gin.Context) {
	var req struct {
		Day int64 `json:"day"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数不全", gin.H{})
		return
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	err = Base.MysqlConn.Delete(&Message.Message{}, "create_time <= ?", startOfDay.AddDate(0, 0, int(-req.Day))).Error
	if err != nil {
		Common.ApiResponse{}.Success(c, err.Error(), gin.H{})
		return
	}
	Common.ApiResponse{}.Success(c, "ok", gin.H{"time": startOfDay.AddDate(0, 0, int(-req.Day))})
}
