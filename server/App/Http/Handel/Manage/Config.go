package Manage

import (
	"github.com/gin-gonic/gin"
	Common2 "server/App/Common"
	"server/App/Http/Request"
	"server/App/Model/Common"
	"server/Base"
)

type Config struct{}

// @summary 配置 获取配置
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Router /manager/config/get [post]
func (Config) Get(c *gin.Context) {
	var config Common.SystemConfig
	Base.MysqlConn.Find(&config)
	Common2.ApiResponse{}.Success(c, "获取成功", gin.H{"config": config})
}

// @summary 配置 修改配置
// @tags 客服系统总后台
// @Param token header string true "认证token"
// @Param system_name query string false  "system_name"
// @Param system_video query string false  "system_video"
// @Param system_version query string false  "system_version"
// @Param pc_download_link query string false  "pc_download_link"
// @Param oss_storage query string false  "oss_storage"
// @Param oss_ali_access_key query string false  "oss_ali_access_key"
// @Param oss_ali_secret_key query string false  "oss_ali_secret_key"
// @Param oss_ali_region query string false  "oss_ali_region"
// @Param oss_ali_domain query string false  "oss_ali_domain"
// @Param oss_tencent_access_key query string false  "oss_tencent_access_key"
// @Param oss_tencent_secret_key query string false  "oss_tencent_secret_key"
// @Param oss_tencent_region query string false  "oss_tencent_region"
// @Param oss_tencent_domain query string false  "oss_tencent_domain"
// @Param ad_default query string false  "ad_default"
// @Param pay_wechat_mch_id query string false  "pay_wechat_mch_id"
// @Param pay_wechat_appid query string false  "pay_wechat_appid"
// @Param pay_wechat_app_key query string false  "pay_wechat_app_key"
// @Param pay_ali_mch_id query string false  "pay_ali_mch_id"
// @Param pay_ali_appid query string false  "pay_ali_appid"
// @Param pay_ali_app_key query string false  "pay_ali_app_key"
// @Router /manager/config/update [post]
func (Config) Update(c *gin.Context) {
	var req Request.SystemConfig
	err := c.ShouldBind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}

	var model Common.SystemConfig
	Base.MysqlConn.Model(&Common.SystemConfig{}).Find(&model)
	if model.Id == 0 {
		Base.MysqlConn.Model(&model).Create(req)
	} else {
		model.RealName = 0
		Base.MysqlConn.Model(&Common.SystemConfig{}).Where("id = ?", model.Id).Updates(req)
	}

	Common2.ApiResponse{}.Success(c, "操作成功", gin.H{"req": req, "model": model})
}
