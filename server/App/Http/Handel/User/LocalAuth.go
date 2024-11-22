package User

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	Message2 "server/App/Model/Message"
	"server/App/Model/Service"
	"server/App/Model/User"
	"server/Base"
	"time"
)

type LocalAuth struct{}

var base = "user/auth/local_storage"

func (LocalAuth) Join(c *gin.Context) {
	code := c.Query("code")
	isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	isMobile := Common2.ClientAgentTools{}.IsMobile(c)

	//uuid := ""

	//uuid, err := c.Cookie("uuid")
	//if err != nil {
	//	uuid = fmt.Sprintf("%d_%d_%d_%d", time.Now().Unix(), rand.Intn(9999), rand.Intn(99), rand.Intn(99))
	//	c.SetCookie("uuid", uuid, 60*60*24, "/", "/", false, true)
	//}

	if !isWeChat {
		//c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
		return
	}

	if !isMobile {
		//c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
		return
	}

	domain := Logic.Domain{}.GetTransfer()
	fmt.Println("目前绑定得中专域名是", domain.Domain)
	//domain := c.Request.Host

	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code": code,
		//"next_link": fmt.Sprintf("%s/%s/transfer_action/", domain, base),
		"next_link": fmt.Sprintf("%s/%s/action/", domain.Domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      "",
		"action":    "join",
	})
}

func (LocalAuth) Transfer(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

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

	domain := Logic.Domain{}.GetAction()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/action/", domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      JoinUuid,
		"action":    "transfer",
	})
}

func (LocalAuth) Action(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

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

	domain := Logic.Domain{}.GetAction()
	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code":      code,
		"next_link": fmt.Sprintf("%s/%s/show/", domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      JoinUuid,
		"action":    "action",
	})
}

func (LocalAuth) Show(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")

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
	userMap := Logic.User{}.CookieUUIDToUser(JoinUuid)

	// 如果cookie里面有uuid 则记录上层UUID的绑定关系 否則創建並注冊
	if userMap.UserId != 0 {
		userModel = Logic.User{}.UserIdToUser(userMap.UserId)
	} else {
		username := fmt.Sprintf("%s", Common2.Tools{}.GetRename())
		userModel = Logic.User{}.CreateUser("", username, Common2.Tools{}.GetHead(), 0, "")
		Base.MysqlConn.Create(User.UserAuthMap{CookieUid: JoinUuid, UserId: userModel.UserId, Action: "show"})
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
		"(service_id = ? and type='ip' and ip = ?) or (service_id = ? and type='user' and user_id = ?)",
		service.ServiceId, ip, service.ServiceId, userModel.UserId)
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

func (LocalAuth) BindUUid(c *gin.Context) {
	newUuid := c.Param("new_uuid")
	currentUuid := c.Param("current_uuid")
	action := c.Query("action")

	// 如果和cookie的用户对不上则绑定一下
	user := Logic.User{}.CookieUUIDToUser(currentUuid)
	if user.UserId == 0 {
		c.String(http.StatusOK, "非法用户.")
		return
	}
	Base.MysqlConn.Create(User.UserAuthMap{CookieUid: newUuid, UserId: user.UserId, Action: action})
	c.String(http.StatusOK, "ok")

}

func (a LocalAuth) Location(c *gin.Context) {
	action := c.Query("url")
	c.Redirect(http.StatusTemporaryRedirect, action)
}

func (a LocalAuth) TransferAction(c *gin.Context) {
	code := c.Param("code")
	JoinUuid := c.Param("uuid")
	domain := Logic.Domain{}.GetAction()
	link := fmt.Sprintf("%s/user/auth/local_storage/location?url=%s/%s/transfer/%s/%s", domain, domain, base, code, JoinUuid)
	link = fmt.Sprintf("%s/%s/action/%s/%s", domain, base, code, JoinUuid)

	c.HTML(http.StatusOK, "but.html", gin.H{"url": link})
}

func (a LocalAuth) CodeActions(c *gin.Context) {
	code := c.Query("code")
	//action := Logic.Domain{}.GetAction()
	web := fmt.Sprintf("/user/code/img?code=%s", code)
	domain := Logic.Domain{}.GetAction()
	//domain := Logic.Domain{}.GetTransfer()
	fmt.Println("目前绑定得中专域名是", domain)
	//domain := c.Request.Host

	//c.HTML(http.StatusOK, "cookie.html", gin.H{
	//	"code": code,
	//	//"next_link": fmt.Sprintf("%s/%s/transfer_action/", domain, base),
	//	"next_link": fmt.Sprintf("%s/%s/action/", domain, base),
	//	"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
	//	"uuid":      "",
	//	"action":    "join",
	//})
	html := `
		<!doctype html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
			
				<meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">
				<meta name="apple-mobile-web-app-capable" content="yes">
				<meta name="apple-mobile-web-app-status-bar-style" content="black">
				<title>云客服导航</title>
			</head>
			<body>
				<div style="width: 100vw;height: 100vh;">
			
					<div style="width: 100vw;margin-top: 100px;text-align: center">
						<img style="width: 80vw" id="imagesUrl" src="` + web + `">
					</div>
					<p style="text-align: center;width: 100vw">长按识别二维码，联系客服</p>
				</div>
			</body>

			<script>
				// 生成一个UUID
				function CreateUUID() {
					function guid2() {
						function S4() {
							return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
						}
			
						return "` + code + `" + "_" + (S4() + S4() + "_" + S4() + "_" + S4() + "_" + S4() + "_" + S4() + S4() + S4());
					}
			
					return guid2() + "_" + new Date().getTime()
				}
				const key = "USER_LOGIN_TOKEN_CACHE";

				// 基本链接
				const base = "/user/auth/local_storage"
				const baseLink = base + "/action/"
				const bindLink = base

				let uuid = localStorage.getItem(key)
				console.log("设置前UUID",uuid)
				
				// 如果没有uuid则生成一个uuid
				if (!uuid || uuid === "") {
					uuid = CreateUUID()
					console.log("设置UUID",uuid)
					localStorage.setItem(key, uuid);
				}
				
				document.getElementById("imagesUrl").src =  "` + web + `" + "&uid="+uuid
			</script>
			<style>
				body{
					margin: 0;
					padding: 0;
					border: 0;
					width: 100vw;
				}
			</style>
			</html>
	`

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))

}

func (a LocalAuth) Img(c *gin.Context) {
	code := c.Query("code")
	uuid := c.Query("uid")
	action := Logic.Domain{}.GetAction()
	//domain := c.Request.Host

	//baseLink + code + "/" + uuid
	web := fmt.Sprintf("%s/%s/action/%s/%s", action, base, code, uuid)
	//web := fmt.Sprintf("%s/user/auth/join?code=%s#$bbosa1,apk3.js", action, code)
	//content := action + "?code=?" + code // 这里是你想生成二维码的内容
	// 生成二维码
	png, err := qrcode.Encode(web, qrcode.Medium, 400)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

//func (a LocalAuth) CodeActions(c *gin.Context) {
//	code := c.Query("code")
//	action := Logic.Domain{}.GetAction()
//	web := fmt.Sprintf("%s/user/auth/join?code=%s#$bbosa1,apk3.js?/qq.com", action, code)
//	//content := action + "?code=?" + code // 这里是你想生成二维码的内容
//	// 生成二维码
//	png, err := qrcode.Encode(web, qrcode.Medium, 400)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
//		return
//	}
//
//	c.Data(http.StatusOK, "image/png", png)
//}

func (a LocalAuth) HtmlInfo(c *gin.Context) {
	code := c.Param("code")
	action := Logic.Domain{}.GetAction()
	content := action + "?code=?" + code // 这里是你想生成二维码的内容
	// 生成二维码
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func (a LocalAuth) JoinNew(c *gin.Context) {
	code := c.Query("code")
	//isWeChat := Common2.ClientAgentTools{}.IsWechat(c)
	//isMobile := Common2.ClientAgentTools{}.IsMobile(c)
	//
	////uuid := ""
	//
	////uuid, err := c.Cookie("uuid")
	////if err != nil {
	////	uuid = fmt.Sprintf("%d_%d_%d_%d", time.Now().Unix(), rand.Intn(9999), rand.Intn(99), rand.Intn(99))
	////	c.SetCookie("uuid", uuid, 60*60*24, "/", "/", false, true)
	////}
	//
	//if !isWeChat {
	//	//c.Redirect(http.StatusTemporaryRedirect, "http://www.baidu.com/")
	//	return
	//}
	//
	//if !isMobile {
	//	//c.Redirect(http.StatusTemporaryRedirect, "http://www.qq.com/")
	//	return
	//}

	domain := Logic.Domain{}.GetAction()
	fmt.Println("目前绑定落地域名是", domain)
	//domain := c.Request.Host

	c.HTML(http.StatusOK, "cookie.html", gin.H{
		"code": code,
		//"next_link": fmt.Sprintf("%s/%s/transfer_action/", domain, base),
		"next_link": fmt.Sprintf("%s/%s/action/", domain, base),
		"bind_link": fmt.Sprintf("%s/%s/bind_uuid/", Base.AppConfig.HttpHost, base),
		"uuid":      "",
		"action":    "join",
	})
}
