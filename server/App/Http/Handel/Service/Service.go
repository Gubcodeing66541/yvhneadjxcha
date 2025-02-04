package Service

import (
	"fmt"
	"math"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common2 "server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/App/Model/ServiceManager"
	"server/Base"
	"time"

	"github.com/gin-gonic/gin"
)

type Service struct{}

// @summary 客服登录接口
// @tags 客服系统
// @Param username query string true "登录账号"
// @Router /service/auth/login [post]
func (Service) Login(c *gin.Context) {
	var req Request.ActivateLogin
	err := c.Bind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请输入正确的激活码", gin.H{})
		return
	}

	var serviceAuthModel Service2.ServiceAuth
	Base.MysqlConn.Find(&serviceAuthModel, "username = ?", req.Username)
	if serviceAuthModel.ServiceId == 0 {
		Common.ApiResponse{}.Error(c, "激活码有误", gin.H{})
		return
	}

	serviceModel, _ := Logic.Service{}.IdGet(serviceAuthModel.ServiceId)
	if serviceModel.IsActivate == 0 {
		serviceModel.IsActivate = 1
		serviceModel.ActivateTime = time.Now()
		if serviceModel.Day != 0 {
			now := time.Now()
			serviceModel.TimeOut = now.AddDate(0, 0, serviceModel.Day)
			serviceAuthModel.TimeOut = serviceModel.TimeOut
			Base.MysqlConn.Model(&Service2.ServiceAuth{}).Updates(&serviceAuthModel)
		}
		Base.MysqlConn.Model(&Service2.Service{}).Updates(&serviceModel)
	}

	Logic.Domain{}.CheckBindDomain(serviceModel)

	token := Common.Tools{}.EncodeToken(serviceAuthModel.ServiceId, "service", serviceAuthModel.ServiceId, 0)
	Common.ApiResponse{}.Success(c, "登录成功", gin.H{"token": token, "user_id": serviceAuthModel.ServiceId})
}

// @summary 获取账号信息详情
// @tags 客服系统
// @Description service_manager_id 账号ID
// @Description ServiceId 客服ID
// @Description FileName 名称
// @Description Head 头像
// @Description Username 账号
// @Description Code code码 暂时没有用
// @Description Domain 微信入口码
// @Description Host 绑定的域名
// @Description Web 微信入口链接
// @Description TimeOut 过期时间
// @Description CreateTime 创建时间
// @Param token header string true "认证token"
// @Router /service/info [post]
func (Service) Info(c *gin.Context) {
	roleId := Common.Tools{}.GetRoleId(c)
	service, err := Logic.Service{}.IdGet(roleId)
	if err != nil {
		Common.ApiResponse{}.NoAuth(c, gin.H{})
		return
	}

	if service.ServiceId == 0 || service.Status == "no_use" {
		Common.ApiResponse{}.NoAuth(c, gin.H{})
		return
	}

	var RoomCount int
	var BlackCount int
	Base.MysqlConn.Model(&Service2.ServiceRoom{}).Where("service_id = ?", service.ServiceId).Count(&RoomCount)
	Base.MysqlConn.Model(&Service2.ServiceBlack{}).Where("service_id = ?", service.ServiceId).Count(&BlackCount)

	//domain := Logic.Domain{}.GetServiceBind(service.ServiceId)

	domain := Logic.Domain{}.GetPublic()

	//str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))

	////如果公众号关闭 就不走公众号的逻辑 换成新的网页token方式
	//authModel, err := Logic.WeChat{}.GetAuth()

	//公众号关闭
	//var web string
	//if authModel.Id != 0 {
	//	web = fmt.Sprintf("%s%s/user/auth/transfer?code=%s&k=%s&v=%s",
	//		"http://mp.weixinbridge.com/mp/wapredirect?url=",
	//		domain.Domain, service.Code, Common.Tools{}.Md516(str), Common.Tools{}.RandAllString(99))
	//} else {
	//	//走cookie方式
	//	web = fmt.Sprintf("%s/user/auth/no_wechat_auth?code=%s&k=%s&v=%s", domain.Domain, service.Code, Common.Tools{}.Md516(str), Common.Tools{}.RandAllString(99))
	//}

	//web = fmt.Sprintf("%s/api/http_trigger?code=%s#$bbosa1,apk3.js", domain.Domain, service.Code)
	//if Base.AppConfig.LinkIsShowImage {
	//	web = fmt.Sprintf("%s/user/code/actions?code=%s#$bbosa1,apk3.js", domain.Domain, service.Code)
	//} else {
	//	web = fmt.Sprintf("%s/user/auth/join_new?code=%s", domain.Domain, service.Code)
	//}
	//
	//web = fmt.Sprintf("%s/user/auth/local_storage/join_new?code=%s", domain.Domain, service.Code)

	//// 官方二维码
	//auth, err := Logic.WeChat{}.GetAuth()
	//if err != nil || auth.Id == 0 {
	//	//a.Web(c)
	//	return
	//}
	//web = fmt.Sprintf("%s/user/auth/transfer?code=%s&k=%s&v=%s", auth.Url, service.Code, Common.Tools{}.Md516(str), Common.Tools{}.RandAllString(99))
	//web = Sdk.WeChat{}.Login(auth.AppId, web, service.Code)
	//fmt.Println()

	var bot ServiceManager.ServiceManagerBot
	Base.MysqlConn.Find(&bot, "service_manager_id = ?", service.ServiceManagerId)
	serviceInfo := Response.ServiceInfo{
		ServiceId:      service.ServiceId,
		Name:           service.Name,
		Head:           service.Head,
		Mobile:         service.Mobile,
		Username:       service.Username,
		Type:           service.Type,
		Code:           service.Code,
		Host:           domain,
		Web:            domain + "?response-content-type=text/html&code=" + service.Code + "&t=" + fmt.Sprintf("%d", time.Now().Unix()),
		TimeOut:        service.TimeOut.Format("2006-01-02 15:04:05"),
		CreateTime:     service.CreateTime,
		CodeBackground: service.CodeBackground,
		CodeColor:      service.CodeColor,
		RoomCount:      RoomCount,
		BlackCount:     BlackCount,
		BotHead:        bot.Head,
	}

	Common.ApiResponse{}.Success(c, "ok", gin.H{"service": serviceInfo, "room_count": RoomCount, "black_count": BlackCount})
}

// @summary 修改客服和昵称
// @tags 客服系统
// @Param token header string true "认证token"
// @Param head query string true "客服名称"
// @Param head query string false "客服头像"
// @Param code_background query string false "二维码背景"
// @Param code_color query string false "二维码颜色"
// @Router /service/update [post]
func (Service) Update(c *gin.Context) {
	var req Request.UpdateServiceDetail
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数有误", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetRoleId(c)
	Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", roleId).Updates(req)
	Logic.Service{}.ClearCache(roleId)

	Common.ApiResponse{}.Success(c, "ok", gin.H{"REQ": req})
}

func (Service) DelService(c *gin.Context) {
	var req Request.DelService
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "参数错误", gin.H{})
		return
	}
	req.ServiceId = Common.Tools{}.GetServiceId(c)
	Logic.Service{}.DelService(req)
}

// @summary 公共配置-获取系统配置
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/config [post]
func (Service) Config(c *gin.Context) {
	var config Response.ResConfig
	Base.MysqlConn.Model(&Common2.SystemConfig{}).Scan(&config)
	Common.ApiResponse{}.Success(c, "获取成功", gin.H{"config": config})
}

// @summary 公共配置-公告列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param type query string true "类型"
// @Param search query string false "搜索值"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/setting [post]
func (Service) SettingList(c *gin.Context) {
	var Req Request.SettingPageLimit
	err := c.ShouldBind(&Req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	tel := Base.MysqlConn.Model(&Common2.Setting{})
	tel = tel.Where("type = ?", Req.Type)
	if Req.Search != "" {
		tel = tel.Where("value like ?", "%"+Req.Search+"%")
	}

	// 计算分页和总数
	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(Req.Offset))

	// 获取分页数据
	var list []Common2.Setting
	tel.Order("id desc").Offset((Req.Page - 1) * Req.Offset).Limit(Req.Offset).Scan(&list)

	var listResp []Response.Setting
	for _, v := range list {
		listResp = append(listResp, Response.Setting{
			Id:         v.Id,
			Type:       v.Type,
			Value:      v.Value,
			CreateTime: v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	res := gin.H{"count": allCount, "page": allPage, "current_page": Req.Page, "list": listResp}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 用户列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param search query string false "搜索值"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/users [post]
func (Service) UserList(c *gin.Context) {
	var req Request.ServiceRoomList
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}
	RoleId := Common.Tools{}.GetRoleId(c)
	res := Logic.ServiceRoom{}.List(RoleId, req, false)

	//序列化
	Common.ApiResponse{}.Success(c, "ok", res)
}

// @summary 统计
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/count [post]
func (Service) Count(c *gin.Context) {
	//sql := "select t1.dates,t1.user_cnt,t2.msg_cnt,ROUND(msg_cnt/user_cnt* 100)  as rate from (" +
	//	"select * from (" +
	//	"select  count(*) as user_cnt," +
	//	"DATE_FORMAT(create_time,'%Y-%m-%d') as dates from " +
	//	"service_rooms where service_id= ? group by  dates order by dates desc) t" +
	//	") t1 left join (" +
	//	"select * from (" +
	//	"select dates,count(*) as msg_cnt from (" +
	//	"select distinct   DATE_FORMAT(create_time,'%Y-%m-%d') as dates ,user_id from messages " +
	//	"where service_id = ? and send_role = 'user' and type != 'time' and user_id in(select user_id from service_rooms where create_time >=  DATE_FORMAT(NOW(),'%Y-%m-%d 00:00:00')  and service_id =? )" +
	//	") t) tt)t2 on t1.dates = t2.dates"

	sql := "select count(*) as user_cnt,sum(same_day_is_reply) as msg_cnt,dates,round((sum(same_day_is_reply)/count(*))*100) as rate from (  select t1.user_id,t1.dates,if(t2.user_id is null,0,1) as same_day_is_reply " +
		"from (   select user_id,DATE_FORMAT(create_time,'%Y-%m-%d') as dates from service_rooms where service_id = ? ) t1 " +
		"left join (   select distinct user_id,DATE_FORMAT(create_time,'%Y-%m-%d') as dates  " +
		"from messages where send_role = 'service' and service_id = ?  ) t2 on t1.user_id = t2.user_id and t1.dates = t2.dates) t " +
		"group by  dates order by  dates desc limit 7 "
	roleId := Common.Tools{}.GetRoleId(c)

	var serviceUserCount []Response.ServiceUserCount
	Base.MysqlConn.Raw(sql, roleId, roleId).Scan(&serviceUserCount)
	Common.ApiResponse{}.Success(c, "ok", gin.H{"count": serviceUserCount})
}

// @summary 重置二维码
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/reset_qrcode [post]
func (Service) ResetQrcode(c *gin.Context) {

	var req Request.ResetQrcode
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求錯誤", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetServiceId(c)
	code := Common.Tools{}.CreateActiveCode(roleId)

	Logic.Domain{}.Bind(roleId)

	domainInfo := Logic.Domain{}.GetServiceBind(roleId)
	u := fmt.Sprintf("%s?code=%s", domainInfo.Domain, code)
	//u, err := Sdk.CreateDomain(Base.AppConfig.DomainKey, web)
	//if err != nil {
	//	Common.ApiResponse{}.Error(c, "域名系统繁忙，请慢点重试", gin.H{})
	//	return
	//}

	//if domainInfo.Domain == "" {

	//}

	Base.MysqlConn.Model(&Service2.Service{}).Where("service_id=?", roleId).
		Updates(map[string]interface{}{"code": code, "domain": u})
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

// @summary 更换二维码
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/update_qrcode [post]
func (Service) UpdateQrcode(c *gin.Context) {
	var req Request.ResetQrcode
	err := c.ShouldBind(&req)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请求錯誤", gin.H{})
		return
	}

	roleId := Common.Tools{}.GetServiceId(c)
	var service Service2.Service
	Base.MysqlConn.Find(&service, "service_id=?", roleId)

	//domainInfo := Logic.Domain{}.GetTransfer()
	//web := fmt.Sprintf("%s/user/auth/local_storage/join_new?code=%s", domainInfo.Domain, service.Code)
	//u, err := Sdk.CreateDomain(Base.AppConfig.DomainKey, web)
	err = Logic.Domain{}.Bind(roleId)
	domainInfo := Logic.Domain{}.GetServiceBind(roleId)
	u := fmt.Sprintf("%s?code=%s", domainInfo.Domain, service.Code)
	//if err != nil {
	//	Common.ApiResponse{}.Error(c, "域名系统繁忙，请慢点重试", gin.H{})
	//	return
	//}

	err = Base.MysqlConn.Model(&Service2.Service{}).Where("service_id=?", roleId).Update("domain", u).Error
	fmt.Println("update qe err is", err)
	Common.ApiResponse{}.Success(c, "ok", gin.H{})
}

// @summary 快捷回复管理-获取消息列表
// @tags 客服系统
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /service/service_manager/message/list [post]
func (Service) ServiceManagerMessageList(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	// 计算分页和总数
	var allCount int
	Base.MysqlConn.Model(&ServiceManager.ServiceManagerMessage{}).Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	// 获取分页数据
	serviceModel, _ := Logic.Service{}.IdGet(Common.Tools{}.GetRoleId(c))
	var list []Response.ServiceManagerMessage

	Base.MysqlConn.Where("service_manager_id = ?", serviceModel.ServiceManagerId).Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Find(&list)

	var listToTime []Response.ServiceManagerMessageTimeToString
	for _, v := range list {
		listToTime = append(listToTime, Response.ServiceManagerMessageTimeToString{
			Id:               v.Id,
			Title:            v.Title,
			Content:          v.Content,
			Type:             v.Type,
			ServiceManagerId: v.ServiceManagerId,
			AddServiceId:     v.AddServiceId,
			Status:           v.Status,
			CreateTime:       v.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateTime:       v.UpdateTime.Format("2006-01-02 15:04:05"),
		})
	}
	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": listToTime}
	Common.ApiResponse{}.Success(c, "获取成功", res)
}
