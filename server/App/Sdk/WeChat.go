package Sdk

import (
	"encoding/json"
	"fmt"
	"server/App/Common"
	"server/App/Http/Response"
)

type WeChat struct {}

func (WeChat) Login(appId string,redirectUri string,state string) string {
	url := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
	return fmt.Sprintf(url,appId,redirectUri,state)
}

func (WeChat) CodeToUserAuthMsg(appId string,secret string,code string) (Response.WeChatAuth,error)  {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",appId,secret,code)
	token := Common.Tools{}.HttpGet(url)
	var weChatAuth Response.WeChatAuth
	err := json.Unmarshal(token,&weChatAuth);if err != nil{
		return weChatAuth,err
	}
	return weChatAuth,nil
}

func (WeChat) GetAccessToken(appId string,appSecret string) []byte {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",appId,appSecret)
	res := Common.Tools{}.HttpGet(url)
	return res
}

func (WeChat) AccessToUserInfo(accessToken string,openId string) []byte  {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",accessToken,openId)
	res := Common.Tools{}.HttpGet(url)
	return res
}