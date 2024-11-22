package Manage

import (
	"math"
	"server/App/Common"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common3 "server/App/Model/Common"
	"server/Base"

	"github.com/gin-gonic/gin"
)

type WeChatAuths struct{}

// @summary 公众号-列表
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param type query string true "类型 push auth"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /manager/wechat_auths/list [post]
func (WeChatAuths) List(c *gin.Context) {
	var Req Request.WeChatType
	err := c.ShouldBind(&Req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	tel := Base.MysqlConn.Model(&Common3.WeChatAuth{})
	if Req.Search != "" {
		tel = tel.Where("name like ?", "%"+Req.Search+"%")
	}

	if Req.Type != "" {
		tel = tel.Where("type = ?", Req.Type)
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(Req.Offset))

	// 获取分页数据
	var list []Common3.WeChatAuth
	tel.Offset((Req.Page - 1) * Req.Offset).Limit(Req.Offset).Find(&list)

	var respList []Response.WeChatAuthResp

	for _, v := range list {
		respList = append(respList, Response.WeChatAuthResp{
			Id:         v.Id,
			Name:       v.Name,
			AppId:      v.AppId,
			AppSecret:  v.AppSecret,
			FileName:   v.FileName,
			FileValue:  v.FileValue,
			Type:       v.Type,
			Url:        v.Url,
			UrlSpare:   v.UrlSpare,
			Status:     v.Status,
			MessageId:  v.MessageId,
			CreateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}
	res := gin.H{"count": allCount, "page": allPage, "current_page": Req.Page, "list": respList}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 公众号-快速启用禁用
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param type query string true "类型 push auth"
// @Param id query int true "id"
// @Param status query string true "状态 puth auth"
// @Router /manager/wechat_auths/enable_disable [post]
func (WeChatAuths) OpenOrClose(c *gin.Context) {
	var req Request.OpenOrClose
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Logic.WeChatAuths{}.OpenOrClose(req.Id, req.Status)
	Common.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 公众号-新增
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param name query string true "名称"
// @Param app_id query string true "appid"
// @Param app_secret query string true "app_secret"
// @Param file_name query string true "文件名.txt"
// @Param file_value query string true "文件值"
// @Param type query string true "类型 auth授权 push推送"
// @Param url query string true "授权地址"
// @Param url_spare query string true "授权备用地址"
// @Router /manager/wechat_auths/create [post]
func (WeChatAuths) Create(c *gin.Context) {
	var req Request.CreateWeChatAuth
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	Logic.WeChatAuths{}.Create(req)
	Common.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 公众号-修改
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query string true "授权的ID"
// @Param name query string true "名称"
// @Param app_id query string true "appid"
// @Param app_secret query string true "app_secret"
// @Param file_name query string true "文件名.txt"
// @Param file_value query string true "文件值"
// @Param type query string true "类型 auth授权 push推送"
// @Param url query string true "授权地址"
// @Param url_spare query string true "授权备用地址"
// @Router /manager/wechat_auths/update [post]
func (WeChatAuths) Update(c *gin.Context) {
	var req Request.UpdateWeChatAuth
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.WeChatAuths{}.Update(req)
	if err2 != nil {
		Common.ApiResponse{}.Error(c, err.Error(), gin.H{})
		return
	}
	Common.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 公众号-删除
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query string true "id"
// @Router /manager/wechat_auths/delete [post]
func (WeChatAuths) Delete(c *gin.Context) {
	var req Request.DeleteWeChatAuth
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	err2 := Logic.WeChatAuths{}.Delete(req.Id)
	if err2 != nil {
		Common.ApiResponse{}.Error(c, err2.Error(), gin.H{})
		return
	}
	Common.ApiResponse{}.Success(c, "操作成功", gin.H{})
}

// @summary 公众号-一键切换url
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param id query string true "id"
// @Param url query string true "url"
// @Param url_spare query string true "url_spare"
// @Router /manager/wechat_auths/switch [post]
func (WeChatAuths) Switch(c *gin.Context) {
	var req Request.WeChatSwitch
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "暂无备用域名，无法切换", gin.H{})
		return
	}
	Base.MysqlConn.Model(&Common3.WeChatAuth{}).Where("id = ?", req.Id).Updates(req)
	Common.ApiResponse{}.Success(c, "操作成功", gin.H{})
}
