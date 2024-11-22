package Common

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type ClientAgentTools struct{}

// 是否是微信检测
func (ClientAgentTools) IsWechat(c *gin.Context) bool {
	return true
	userAgent := c.GetHeader("User-Agent")
	if strings.Index(userAgent, "MicroMessenger") != -1 {
		return true
	}
	c.HTML(403, "nginx.html", "")
	c.Abort()
	return false
}

func (ClientAgentTools) IsMobile(c *gin.Context) bool {
	//return true
	isMobile := false

	userAgent := c.GetHeader("User-Agent")
	if len(userAgent) == 0 {
		isMobile = false
	}

	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	if !isMobile {
		c.HTML(403, "nginx.html", "")
		c.Abort()
	}
	return isMobile
}

func (ClientAgentTools) GetDrive(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	if strings.Index(userAgent, "Windows") != -1 {
		return "Windows"
	}
	if strings.Index(userAgent, "Android") != -1 {
		return "Android"
	}
	if strings.Index(userAgent, "iPhone") != -1 {
		return "iPhone"
	}
	if strings.Index(userAgent, "iPod") != -1 {
		return "iPod"
	}
	if strings.Index(userAgent, "iPad") != -1 {
		return "iPad"
	}
	if strings.Index(userAgent, "Windows Phone") != -1 {
		return "Windows Phone"
	}
	if strings.Index(userAgent, "MQQBrowser") != -1 {
		return "QQ浏览器"
	}
	if strings.Index(userAgent, "iPhone") != -1 {
		return "iPhone"
	}
	return "未知"
}
