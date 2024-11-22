package Logic

import (
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Model/Common"
	"server/Base"
	"time"
)

// 微信管理
type WeChatAuths struct {
}

// 微信公众号列表
func (WeChatAuths) List(req Request.WeChatList) Response.ResultList {

	var result Response.ResultList
	page := req.Page
	result.Page = page

	var weChatAuth []Common.WeChatAuth
	Base.MysqlConn.Find(&weChatAuth)
	result.Page.Count = len(weChatAuth)

	var respWeChatAuthList []Common.WeChatAuth
	if req.Type == "" {
		Base.MysqlConn.Raw(" select * from we_chat_auths  limit ?,? ", (req.Page.CurrentPage-1)*req.Page.CurrentSize, req.Page.CurrentSize).Scan(&respWeChatAuthList)
	} else {
		Base.MysqlConn.Raw(" select * from we_chat_auths where type = ? limit ?,? ", req.Type, (req.Page.CurrentPage-1)*req.Page.CurrentSize, req.Page.CurrentSize).Scan(&respWeChatAuthList)
	}
	result.Data = respWeChatAuthList
	return result
}

func (WeChatAuths) OpenOrClose(id int, status string) {
	var weChatAuths Common.WeChatAuth
	Base.MysqlConn.Model(&weChatAuths).Where("id = ? ", id).Update("status", status)
	WeChat{}.ClearCache()
}

func (WeChatAuths) Create(req Request.CreateWeChatAuth) {
	Base.MysqlConn.Create(&Common.WeChatAuth{
		Name: req.Name, AppId: req.AppId, AppSecret: req.AppSecret, FileName: req.FileName, FileValue: req.FileValue,
		Type: req.Type, Url: req.Url, UrlSpare: req.UrlSpare, Status: "un_enable", CreateTime: time.Now(), MessageId: req.MessageId,
	})
	WeChat{}.ClearCache()
}

func (WeChatAuths) Update(req Request.UpdateWeChatAuth) error {
	var weChatAuth Common.WeChatAuth
	Base.MysqlConn.Model(&weChatAuth).Updates(&req)
	WeChat{}.ClearCache()
	return nil
}

func (WeChatAuths) Delete(id int) error {
	Base.MysqlConn.Delete(&Common.WeChatAuth{}, "id = ?", id)
	WeChat{}.ClearCache()
	return nil
}

// 获取配置
func (WeChatAuths) getById(id int) {
	var weChatAuth Common.WeChatAuth
	Base.MysqlConn.Find(&weChatAuth, "id = ?", id)
}
