package Manage

import (
	"github.com/gin-gonic/gin"
	"math"
	Common2 "server/App/Common"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common3 "server/App/Model/Common"
	"server/Base"
	"time"
)

type Setting struct{}

// @summary 设定公告敏感词-列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param type query string true "类型"
// @Param search query string false "搜索值"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/setting/list [post]
func (Setting) List(c *gin.Context) {
	var Req Request.SettingPageLimit
	err := c.ShouldBind(&Req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	tel := Base.MysqlConn.Model(&Common3.Setting{})
	tel = tel.Where("type = ?", Req.Type)
	if Req.Search != "" {
		tel = tel.Where("value like ?", "%"+Req.Search+"%")
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(Req.Offset))

	// 获取分页数据
	var list []Common3.Setting
	tel.Offset((Req.Page - 1) * Req.Offset).Limit(Req.Offset).Find(&list)

	var SettingResponse []Response.SettingResponse

	for _, setting := range list {
		SettingResponse = append(SettingResponse, Response.SettingResponse{
			Id:         setting.Id,
			Type:       setting.Type,
			Value:      setting.Value,
			CreateTime: setting.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	res := gin.H{"count": allCount, "page": allPage, "current_page": Req.Page, "list": SettingResponse}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 设定公告敏感词-创建
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param type query string true "类型"
// @Param value query string true "值"
// @Router /manager/setting/create [post]
func (Setting) Create(c *gin.Context) {
	var req Request.CreateSetting
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}
	req.CreateTime = time.Now()
	create := Common3.Setting{Type: req.Type, Value: req.Value, CreateTime: time.Now()}
	Base.MysqlConn.Create(&create)
	Common2.ApiResponse{}.Success(c, "添加成功", gin.H{})
}

// @summary 设定公告敏感词-修改
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query int true "id"
// @Param type query string true "类型 keyword 关键字 notice  公告"
// @Param value query string true "值"
// @Router /manager/setting/update [post]
func (Setting) Update(c *gin.Context) {
	var req Request.CreateSettingUpdate
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请填写完整的数据信息", gin.H{})
		return
	}
	req.CreateTime = time.Now()
	update := Common3.Setting{Type: req.Type, Value: req.Value, CreateTime: time.Time{}}
	Base.MysqlConn.Where("id = ?", req.Id).Model(&Common3.Setting{}).Updates(&update)
	Common2.ApiResponse{}.Success(c, "修改成功", gin.H{})
}

// @summary 设定公告敏感词-删除
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query int true "id"
// @Router /manager/setting/delete [post]
func (Setting) Delete(c *gin.Context) {
	var req Request.Id
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请选择需要删除的数据", gin.H{})
		return
	}
	Base.MysqlConn.Delete(&Common3.Setting{}, "id = ?", req.Id)
	Common2.ApiResponse{}.Success(c, "删除成功", gin.H{})
}
