package Common

import (
	"errors"
	"github.com/gin-gonic/gin"
	"math"
	"server/App/Http/Request"
	"server/Base"
)

type DbHelper struct{}

// 分页器
func (DbHelper) PageTable(model interface{}, modelList interface{}, c *gin.Context) (gin.H, error) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		return nil, errors.New("请提交完整的分页参数")
	}

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Model(&model).Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	ServiceManagerId := Tools{}.GetRoleId(c)
	Base.MysqlConn.Where("service_manager_id = ?", ServiceManagerId).Offset(pageReq.Page - 1*pageReq.Offset).Limit(pageReq.Offset).Find(&modelList)
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": modelList}
	return res, nil
}
