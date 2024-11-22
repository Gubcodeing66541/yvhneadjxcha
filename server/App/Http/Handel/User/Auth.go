package User

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Http/Constant"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	"server/App/Http/Tel"
	Message2 "server/App/Model/Message"
	"server/App/Model/Service"
	"server/App/Model/ServiceManager"
	"server/App/Model/User"
	"server/App/Sdk"
	"server/Base"
	"time"
)

type Auth struct{}

func (Auth) getJoinKey() string {
	return "USER-LOGIN:KEY"
}

// web 万能
func (a Auth) Web(c *gin.Context) {

	code := c.Query("code")
	isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)

	if !isWeChat {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
		return
	}

	if !isMobile {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return
	}

	// 检测如果是微信跳转微信授权逻辑
	service, err := Logic.Service{}.Get(code)
	if err != nil || service.Id == 0 {
		Common2.ApiResponse{}.Error(c, "无法查询的客服信息", gin.H{})
		return
	}

	// 检测是否登陆过，如果登陆过直接跳 否则创建新用户
	actionkey, err := c.Cookie(code)
	if err != nil && actionkey != "" {
		token, err := c.Cookie("token")
		if err == nil && token != "" {
			Common2.RedisTools{}.SetString(actionkey, token)
			action := Logic.Domain{}.GetAction()
			link := fmt.Sprintf("%s/%s", action, actionkey)

			c.Redirect(http.StatusTemporaryRedirect, link)
			return
		}
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 获取随机图片
	files, _ := ioutil.ReadDir("D:\\code\\new_chat\\latest-customer-service-system\\server\\static\\static\\head\\")
	rand.Seed(time.Now().UnixNano())

	username := fmt.Sprintf("%s", Common2.Tools{}.CreateUserName(service.Id))
	userModel := Logic.User{}.CreateUser("", username, Base.AppConfig.HttpHost+"/static/static/head/"+files[rand.Intn(len(files))].Name(), 0, "")
	_ = Logic.ServiceRoom{}.Get(userModel, service.ServiceId, c.ClientIP(), Common2.ClientAgentTools{}.GetDrive(c))

	token := Common2.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0)
	action := Logic.Domain{}.GetAction()
	if action == "" {
		Common2.ApiResponse{}.Error(c, "无法找到action", gin.H{})
		return
	}

	key := fmt.Sprintf("%s%d", Common2.Tools{}.Md516(token), userModel.UserId)
	Common2.RedisTools{}.SetString(key, token)
	c.SetCookie(code, key, 60*60*24*90, "/", "", false, true)
	c.SetCookie("token", token, 60*60*24*90, "/", "", false, true)
	link := fmt.Sprintf("%s/user/auth/action/%s", action, key)
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func (a Auth) Error(c *gin.Context) {
	msg := c.Param("msg")
	Common2.ApiResponse{}.Error(c, msg, gin.H{})
}

// LoginCooKieName cookie授權名字
var LoginCooKieName = "SYSTEM_NO_WECAHRT_LOGIN_UUID"

// NoWeChatAuth 不走微信授权
func (a Auth) NoWeChatAuth(c *gin.Context) {
	code := c.Query("code")

	UUId, err := c.Cookie(LoginCooKieName)
	if err != nil || UUId == "" {
		UUId = Common2.Tools{}.CreateUUID(code)
		outTIme := 60 * 60 * 24 * 30 * 12 * 10
		c.SetCookie(LoginCooKieName, UUId, outTIme, "/", "/", false, true)
	}

	domain := Logic.Domain{}.GetTransfer()
	link := fmt.Sprintf("%s/user/auth/no_wechat_transfer/%s/%s", domain.Domain, code, UUId)

	//c.String(200, fmt.Sprintf("code:%s \r\n cookie:%s \r\n url %s", code, UUId, link))
	//return
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func (a Auth) NoWechatTransfer(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	// 如果UUID不存在则获取上层连接发送出来的uuid并存入cookie
	CookieUUId, err := c.Cookie(LoginCooKieName)
	if err != nil || CookieUUId == "" {
		outTIme := 60 * 60 * 24 * 30 * 12 * 10
		c.SetCookie(LoginCooKieName, JoinUuid, outTIme, "/", "/", false, true)

		domain := Logic.Domain{}.GetAction()
		link := fmt.Sprintf("%s/user/auth/no_wechat_action/%s/%s", domain, code, JoinUuid)
		c.HTML(http.StatusOK, "but.html", gin.H{"url": link})

		//c.Redirect(http.StatusTemporaryRedirect, link)
		return
	}

	// 如果和cookie的用户对不上则绑定一下
	if CookieUUId != JoinUuid {
		user := Logic.User{}.CookieUUIDToUser(CookieUUId)
		if user.UserId == 0 {
			c.String(http.StatusOK, "非法用户.")
			return
		}
		Base.MysqlConn.Create(User.UserAuthMap{CookieUid: JoinUuid, UserId: user.UserId, Action: "NoWechatTransfer!="})
	}

	// 如果cookie里面有uuid切和上层不相同 则记录上层UUID的绑定关系
	domain := Logic.Domain{}.GetAction()
	link := fmt.Sprintf("%s/user/auth/no_wechat_action/%s/%s", domain, code, CookieUUId)
	c.HTML(http.StatusOK, "but.html", gin.H{"url": link})

	//c.Redirect(http.StatusTemporaryRedirect, link)
}

// NoWechatAction 落地
func (a Auth) NoWechatAction(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

	service, err := Logic.Service{}.Get(code)
	if err != nil {
		c.String(http.StatusOK, "未知客服")
		return
	}

	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	// 准备绑定的用户
	var userModel User.User

	// 如果UUID不存在则注册用户并创建cookie、
	CookieUUId, err := c.Cookie(LoginCooKieName)
	if err != nil || CookieUUId == "" {
		outTIme := 60 * 60 * 24 * 30 * 12 * 10
		c.SetCookie(LoginCooKieName, JoinUuid, outTIme, "/", "/", false, true)
		username := Common2.Tools{}.GetRename()
		userModel = Logic.User{}.CreateUser("", username, Common2.Tools{}.GetHead(), 0, "")
		Base.MysqlConn.Create(User.UserAuthMap{CookieUid: JoinUuid, UserId: userModel.UserId, Action: "NoWechatAction=''"})

	} else {
		// 如果cookie里面有uuid 则记录上层UUID的绑定关系
		userModel = Logic.User{}.CookieUUIDToUser(CookieUUId)
		if CookieUUId != JoinUuid {
			Base.MysqlConn.Create(User.UserAuthMap{CookieUid: JoinUuid, UserId: userModel.UserId, Action: "NoWechatAction!="})
		}
	}

	_ = Logic.ServiceRoom{}.Get(userModel, service.ServiceId, c.ClientIP(), Common2.ClientAgentTools{}.GetDrive(c))
	token := Common2.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0)
	action := Logic.Domain{}.GetAction()
	if action == "" {
		Common2.ApiResponse{}.Error(c, "无法找到action", gin.H{})
		return
	}

	ip := c.ClientIP()
	var black Service.ServiceBlack
	Base.MysqlConn.Find(
		&black,
		"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?) or (service_manager_id = ? and user_id = ?)",
		service.ServiceId, ip, service.ServiceId, userModel.UserId, service.ServiceManagerId, userModel.UserId)
	if black.Id != 0 {
		c.String(200, "..")
		return
	}

	// 上线更新update时间和未读
	update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
	Base.MysqlConn.Model(&Service.ServiceRoom{}).
		Where("user_id = ? and service_id = ? ", userModel.UserId, service.ServiceId).
		Updates(update)

	Base.MysqlConn.Model(&Request.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"ip": c.ClientIP()})

	//Base.MysqlConn.Create(&User.UserLoginLog{
	//	UserId:     userModel.UserId,
	//	ServiceId:  service.ServiceId,
	//	Ip:         c.ClientIP(),
	//	Addr:       "",
	//	CreateTime: time.Now(),
	//})

	// 所有消息已读
	Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
		service.ServiceId, userModel.UserId).Updates(
		gin.H{"is_read": 1})

	key := fmt.Sprintf("%s%d", Common2.Tools{}.Md516(token), userModel.UserId)
	Common2.RedisTools{}.SetString(key, token)
	c.SetCookie(code, key, 60*60*24*90, "/", "", false, true)
	c.SetCookie("token", token, 60*60*24*90, "/", "", false, true)
	link := fmt.Sprintf("%s/user/auth/action/%s", action, key)
	c.Redirect(http.StatusTemporaryRedirect, link)
}

func (a Auth) Transfer(c *gin.Context) {
	code := c.Query("code")

	isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)

	if !isWeChat {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
		return
	}

	if !isMobile {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return
	}

	domain := Logic.Domain{}.GetTransfer()

	str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))

	link := fmt.Sprintf("%s%s%s&k=%s&v=%s#%d", domain.Domain, "/user/auth/join?code=", code, Common2.Tools{}.Md516(str), Common2.Tools{}.RandAllString(99), rand.Intn(999999))

	c.Redirect(http.StatusTemporaryRedirect, link)

}

func (a Auth) Join(c *gin.Context) {
	code := c.Query("code")

	isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)

	if !isWeChat {
		//c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
		return
	}

	if !isMobile {
		//c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return
	}

	auth, err := Logic.WeChat{}.GetAuth()
	if err != nil || auth.Id == 0 {
		LocalAuth{}.Join(c)
		return
	}

	service, err := Logic.Service{}.Get(code)
	if err != nil || service.Id == 0 {
		Common2.ApiResponse{}.Error(c, "无法查询的客服信息", gin.H{})
		return
	}

	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	link := fmt.Sprintf("%s%s", auth.Url, "/user/auth/we_chat_auth")
	link = Sdk.WeChat{}.Login(auth.AppId, link, code)

	//正常跳轉
	c.HTML(http.StatusOK, "but.html", gin.H{"url": link})
}

func (a Auth) WebJoin(c *gin.Context) {
	auth, err := Logic.WeChat{}.GetAuth()
	if err != nil || auth.Id == 0 {
		//a.Web(c)
		return
	}

	code := c.Query("code")
	service, err := Logic.Service{}.Get(code)
	if err != nil || service.Id == 0 {
		Common2.ApiResponse{}.Error(c, "无法查询的客服信息", gin.H{})
		return
	}

	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	link := fmt.Sprintf("%s%s", auth.Url, "/user/auth/no_wechat_auth")
	link = Sdk.WeChat{}.Login(auth.AppId, link, code)

	//正常跳轉
	c.HTML(http.StatusOK, "but.html", gin.H{"url": link})
}

// 微信授权
func (Auth) WeChatAuth(c *gin.Context) {
	var req Request.WeChatAuth
	err := c.Bind(&req)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "非法的登录", gin.H{})
		return
	}

	service, err := Logic.Service{}.Get(req.State)
	if err != nil || service.Id == 0 {
		Common2.ApiResponse{}.Error(c, "客服不存在", gin.H{})
		return
	}

	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		Common2.ApiResponse{}.Error(c, "账号已过期", gin.H{})
		return
	}

	//微信授权是否关闭
	authModel, err := Logic.WeChat{}.GetAuth()
	if err != nil || authModel.Id == 0 {
		Common2.ApiResponse{}.Error(c, "无授权信息", gin.H{"authModel": authModel, "req": req})
		return
	}

	var weChatAuth Response.WeChatAuth
	weChatAuth, err = Sdk.WeChat{}.CodeToUserAuthMsg(authModel.AppId, authModel.AppSecret, req.Code)
	if err != nil || weChatAuth.Openid == "" {
		Common2.ApiResponse{}.Error(c, "授权失败", gin.H{})
		return
	}

	// 沒有Uid就用openid
	if weChatAuth.Unionid == "" {
		weChatAuth.Unionid = weChatAuth.Openid
	}

	userModel := Logic.User{}.UnionIdToUser(weChatAuth.Unionid)
	if userModel.UserId == 0 {
		msg := Sdk.WeChat{}.AccessToUserInfo(weChatAuth.AccessToken, weChatAuth.Openid)
		var userReq Response.WeChatUserInfo
		err = json.Unmarshal(msg, &userReq)
		if err != nil {
			Common2.ApiResponse{}.Error(c, "解析用户信息失败", gin.H{"data": msg})
			return
		}
		userModel = Logic.User{}.CreateUser(weChatAuth.Openid, userReq.Nickname, userReq.HeadImgUrl, userReq.Sex, weChatAuth.Unionid)
	}

	rooms := Logic.ServiceRoom{}.Get(userModel, service.ServiceId, c.ClientIP(), Common2.ClientAgentTools{}.GetDrive(c))
	if rooms.IsBlack == 1 {
		c.String(200, ".")
		return
	}

	{
		var serviceBlack Service.ServiceBlack
		Base.MysqlConn.Find(&serviceBlack, "service_manager_id =? and ip= ?", service.ServiceManagerId, c.ClientIP())
		if serviceBlack.Id != 0 {
			c.String(200, ".")
			return
		}

		Base.MysqlConn.Find(&serviceBlack, "service_manager_id =? and user_id= ?", service.ServiceManagerId, rooms.UserId)
		if serviceBlack.Id != 0 {
			c.String(200, ".")
			return
		}
	}

	token := Common2.Tools{}.EncodeToken(userModel.UserId, "user", service.ServiceId, 0)
	action := Logic.Domain{}.GetAction()
	if action == "" {
		Common2.ApiResponse{}.Error(c, "无法找到action", gin.H{})
		return
	}

	key := fmt.Sprintf("%s%d", Common2.Tools{}.Md516(token), userModel.UserId)
	Common2.RedisTools{}.SetString(key, token)
	link := fmt.Sprintf("%s/user/auth/action/%s", action, key)
	c.Redirect(http.StatusTemporaryRedirect, link)
}

// 落地
func (Auth) Action(c *gin.Context) {
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)
	if !isMobile {
		c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return
	}
	// 通过如果token的code存在则检测
	tokenName := c.Param("token")

	//a := c.Query("a")
	//if a == "" {
	//	domain := Logic.Domain{}.GetAction()
	//	url := fmt.Sprintf("%s/user/auth/action/%s?a=44444", domain, token)
	//	fmt.Println("aaa---------------------------", a)
	//	//c.String(200, url)
	//	//return
	//	c.Redirect(http.StatusTemporaryRedirect, url)
	//	return
	//}
	// token 不存在用CookieToken
	if tokenName == "" {
		Common2.ApiResponse{}.Error(c, "未授权", gin.H{})
		return
	}

	token := Common2.RedisTools{}.GetString(tokenName)
	Common2.RedisTools{}.SetString(tokenName, "")

	var userAuthToken Constant.UserAuthToken
	err := Common2.Tools{}.DecodeToken(token, &userAuthToken)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"token": token})
		return
	}

	ip := c.ClientIP()
	var black Service.ServiceBlack
	Base.MysqlConn.Find(
		&black,
		"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?)",
		userAuthToken.ServiceId, ip, userAuthToken.ServiceId, userAuthToken.RoleId)
	if black.Id != 0 {
		c.String(200, "..")
		return
	}

	// 上线更新update时间和未读
	update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
	Base.MysqlConn.Model(&Service.ServiceRoom{}).
		Where("user_id = ? and service_id = ? ", userAuthToken.RoleId, userAuthToken.ServiceId).
		Updates(update)

	Base.MysqlConn.Model(&Request.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
		userAuthToken.ServiceId, userAuthToken.RoleId).Updates(
		gin.H{"ip": c.ClientIP()})

	Base.MysqlConn.Create(&User.UserLoginLog{
		UserId:     userAuthToken.RoleId,
		ServiceId:  userAuthToken.ServiceId,
		Ip:         c.ClientIP(),
		Addr:       "",
		CreateTime: time.Now(),
	})

	// 所有消息已读
	Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
		userAuthToken.ServiceId, userAuthToken.RoleId).Updates(
		gin.H{"is_read": 1})

	c.HTML(http.StatusOK, "index.html", gin.H{"token": token})
	return
}

// @summary 用户token字符串换取真实token
// @tags 用户端
// @Param token header string true "认证token"
// @Router /user/auth/token [post]
func (Auth) Token(c *gin.Context) {
	var req Request.Token
	err := c.ShouldBind(&req)

	if err != nil {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"token": ""})
		return
	}

	token := Common2.RedisTools{}.GetString(req.Token)
	Common2.ApiResponse{}.Success(c, "解析成功", gin.H{"token": token})
}

// @summary 获取客服基本信息
// @tags 用户端
// @Param token query string true "token"
// @Router /user/auth/info [post]
func (Auth) Info(c *gin.Context) {
	serviceId := Common2.Tools{}.GetServiceId(c)

	// 获取客服信息
	service, err := Logic.Service{}.IdGet(serviceId)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "无法获取到客服信息", gin.H{})
		return
	}

	// 获取自己的信息
	users := Logic.User{}.UserIdToUser(Common2.Tools{}.GetRoleId(c))

	//ip := c.ClientIP()
	//var black Service.ServiceBlack
	//Base.MysqlConn.Find(
	//	&black,
	//	"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?)",
	//	service.ServiceId, ip, service.ServiceId, users.UserId)
	//if black.Id != 0 {
	//	Common2.ApiResponse{}.NoAuth(c, gin.H{})
	//	return
	//}
	//
	//// 上线更新update时间和未读
	//update := gin.H{"update_time": time.Now(), "user_no_read": 0, "late_user_read_id": 0, "is_delete": 0, "late_ip": c.ClientIP()}
	//Base.MysqlConn.Model(&Service.ServiceRoom{}).
	//	Where("user_id = ? and service_id = ? ", users.UserId, service.ServiceId).
	//	Updates(update)

	//Base.MysqlConn.Model(&Request.ServiceRoomDetail{}).Where("service_id = ? and user_id = ?",
	//	service.ServiceId, users.UserId).Updates(
	//	gin.H{"ip": c.ClientIP()})
	//
	//Base.MysqlConn.Create(&User.UserLoginLog{
	//	UserId:     users.UserId,
	//	ServiceId:  service.ServiceId,
	//	Ip:         c.ClientIP(),
	//	Addr:       "",
	//	CreateTime: time.Now(),
	//})
	//
	//// 所有消息已读
	//Base.MysqlConn.Model(Message2.Message{}).Where("service_id = ? and user_id = ? and is_read = 0",
	//	service.ServiceId, users.UserId).Updates(
	//	gin.H{"is_read": 1})

	var menu []Service.ServiceMenuSetting
	Base.MysqlConn.Order("sort asc").Find(&menu, "service_id = ?", serviceId)

	var notice Service.ServiceNoticeSetting
	Base.MysqlConn.Find(&notice, "service_id = ?", serviceId)

	var bot ServiceManager.ServiceManagerBot
	Base.MysqlConn.Find(&bot, "service_manager_id = ?", service.ServiceManagerId)

	domain := Logic.Domain{}.GetServiceBind(service.ServiceId)
	str := fmt.Sprintf("MEMBER-%s-time:-rand-%d", time.Now(), rand.Intn(99999))
	web := fmt.Sprintf("%s/user/auth/join?code=%s&k=%s&v=%s", domain.Domain, service.Code, Common2.Tools{}.Md516(str), Common2.Tools{}.RandAllString(99))

	Common2.ApiResponse{}.Success(c, "解析成功", gin.H{
		"notice": notice,
		"menu":   menu,
		"users":  users,
		"service": Response.ServiceInfo{
			ServiceId:      service.ServiceId,
			Name:           service.Name,
			Head:           service.Head,
			Mobile:         service.Mobile,
			Username:       service.Username,
			Type:           service.Type,
			Code:           service.Code,
			Host:           "",
			Web:            web,
			TimeOut:        service.TimeOut.Format("2006-01-02 15:04:05"),
			CreateTime:     service.CreateTime,
			CodeBackground: service.CodeBackground,
			CodeColor:      service.CodeColor,
			RoomCount:      0,
			BlackCount:     0,
			BotHead:        bot.Head,
		},
		"bot_head": bot.Head,
		"token":    Common2.Tools{}.GetCookieToken(c),
	})
}

// @summary 消息列表
// @tags 用户端
// @Param token header string true "认证token"
// @Param page query string true "指定页"
// @Param offset query string true "分页数量"
// @Router /user/message/list [post]
func (Auth) List(c *gin.Context) {
	var pageReq Request.PageLimit
	err := c.ShouldBind(&pageReq)
	if err != nil {
		Common2.ApiResponse{}.Error(c, "请提交完整的分页参数", gin.H{})
		return
	}

	roleId := Common2.Tools{}.GetRoleId(c)
	tel := Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ? and type!='time' and type!='remind' and is_del=0",
		Common2.Tools{}.GetServiceId(c), roleId).Order("id desc")

	var allCount int
	tel.Count(&allCount)
	allPage := math.Ceil(float64(allCount) / float64(pageReq.Offset))

	var list []Tel.Message
	tel.Offset((pageReq.Page - 1) * pageReq.Offset).Limit(pageReq.Offset).Scan(&list)

	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	res := gin.H{"count": allCount, "page": allPage, "current_page": pageReq.Page, "list": list}
	Common2.ApiResponse{}.Success(c, "获取成功", res)
}

// @summary 发送消息给客服
// @tags 用户端
// @Param token header string true "认证token"
// @Param type query string true "消息类型"
// @Param content query string true "消息内容"
// @Router /user/auth/send [post]
func (Auth) Send(c *gin.Context) {
	var req Request.UserSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	UserId := Common2.Tools{}.GetRoleId(c)
	ServiceId := Common2.Tools{}.GetServiceId(c)
	RoomId := Common2.Tools{}.ConvertUserMessageRoomId(ServiceId, UserId)
	serviceIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "service", ServiceId))
	model := &Message2.Message{
		RoomId: RoomId, From: UserId, To: ServiceId, Type: req.Type, Content: req.Content, ServiceId: ServiceId,
		SendRole: "user", CreateTime: time.Now(), IsRead: serviceIsOnline, UserId: UserId}

	if req.Type != "time" {
		model.IsRead = 1
	}

	Base.MysqlConn.Create(&model)

	if model.Id == 0 {
		Common2.ApiResponse{}.Error(c, "消息发送失败", gin.H{})
		return
	}

	sendMsg := Tel.SocketMessage{
		Id:     model.Id,
		RoomId: RoomId, From: UserId, To: ServiceId, Type: req.Type, Content: req.Content, ServiceId: ServiceId,
		SendRole: "user", CreateTime: time.Now().Format("2006-01-02 15:04:05"), IsRead: serviceIsOnline, UserId: UserId,
	}

	if req.Type != "time" {
		LateMsg, LateType := req.Content, req.Type
		update := gin.H{
			"late_role": "user", "late_msg": LateMsg, "update_time": time.Now(), "late_id": model.Id, "is_delete": 0,
			"service_no_read": 0, "user_no_read": 0, "late_type": LateType, "LateUserReadId": model.Id,
		}
		bindUserId := Base.WebsocketHub.GetBindUser(Common2.Tools{}.GetServiceWebSocketId(ServiceId))
		if serviceIsOnline == 0 || bindUserId != UserId {
			update["service_no_read"] = gorm.Expr("service_no_read + ?", 1)
		}
		Base.MysqlConn.Model(&Service.ServiceRoom{}).Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(update)
	}

	// 给客服和用户推送
	Common2.ApiResponse{}.SendMsgToService(ServiceId, "message", sendMsg)
	Common2.ApiResponse{}.SendMsgToUser(UserId, "message", sendMsg)

	// 处理离线消息
	Logic.User{}.HandelLeaveMessage(c, ServiceId, UserId)

	// 机器人
	Logic.User{}.HandelBotMessage(c, ServiceId, UserId, req)

	// 更新已读
	go func(ServiceId int, UserId int) {
		Base.MysqlConn.Model(&Message2.Message{}).Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(gin.H{"is_read": 1})
	}(ServiceId, UserId)

	// ok信息
	Common2.ApiResponse{}.Success(c, "消息发送成功.", gin.H{})
}

func (Auth) Button(c *gin.Context) {
	var req Request.UserSendMessage
	err := c.ShouldBind(&req)
	if err != nil || req.Content == "" {
		Common2.ApiResponse{}.Error(c, "请输入需要发送的消息.", gin.H{"erq": req})
		return
	}

	UserId := Common2.Tools{}.GetRoleId(c)
	ServiceId := Common2.Tools{}.GetServiceId(c)
	RoomId := Common2.Tools{}.ConvertUserMessageRoomId(ServiceId, UserId)
	serviceIsOnline := Base.WebsocketHub.UserIdIsOnline(fmt.Sprintf("%s:%d", "service", ServiceId))
	model := &Message2.Message{
		RoomId: RoomId, From: UserId, To: ServiceId, Type: req.Type, Content: req.Content, ServiceId: ServiceId,
		SendRole: "user", CreateTime: time.Now(), IsRead: serviceIsOnline, UserId: UserId}
	Base.MysqlConn.Create(&model)
	if model.Id == 0 {
		Common2.ApiResponse{}.Error(c, "消息发送失败", gin.H{})
		return
	}

	//更新最后一条信息
	LateMsg, LateType := req.Content, req.Type
	update := gin.H{
		"late_role": "user", "late_msg": LateMsg, "update_time": time.Now(), "late_id": model.Id, "is_delete": 0,
		"service_no_read": 0, "user_no_read": 0, "late_type": LateType, "LateUserReadId": model.Id,
	}
	bindUserId := Base.WebsocketHub.GetBindUser(Common2.Tools{}.GetServiceWebSocketId(ServiceId))
	if serviceIsOnline == 0 || bindUserId != UserId {
		update["service_no_read"] = gorm.Expr("service_no_read + ?", 1)
	}
	Base.MysqlConn.Model(&Service.ServiceRoom{}).Where("service_id = ? and user_id = ?", ServiceId, UserId).Updates(update)

	// 给客服和用户推送
	Common2.ApiResponse{}.SendMsgToService(ServiceId, "message", model)
	Common2.ApiResponse{}.SendMsgToUser(UserId, "message", model)

	//按钮的自动回复
	err = Logic.Message{}.SendToUser(ServiceId, UserId, "text", req.ServiceContent, true)
	if err != nil {
		Common2.ApiResponse{}.Success(c, "消息回复异常.", gin.H{"req": req})
	}

	// 处理离线消息
	Logic.User{}.HandelLeaveMessage(c, ServiceId, UserId)

	// ok信息
	Common2.ApiResponse{}.Success(c, "消息已发送.", gin.H{"req": req})

}

func (a Auth) Test(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")
	c.String(http.StatusOK, fmt.Sprintf("code:%s \r\n uuid %s", code, JoinUuid))
}
