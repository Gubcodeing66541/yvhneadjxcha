package Sdk

import (
	"encoding/json"
	"fmt"
	"math/rand"
	url2 "net/url"
	"server/App/Common"
)

// 195-196-197这三个接口，设置的独享专用接口
//var key = "DNwJLJGfyRZoEBrn41g95QlZPYTtkDbk"

var types = []int64{195, 196, 197}

type Response struct {
	Code  string `json:"code"`
	Id    string `json:"id"`
	Vip   string `json:"vip"`
	Short string `json:"short"`
	Long  string `json:"long"`
	Time  string `json:"time"`
}

func CreateDomain(key, domain string) (url string, err error) {
	typeN := rand.Intn(2)
	api := fmt.Sprintf("https://cdn.yiyai.top/api/?method=add&type=1&key=%s&domain=%d&vip=1&url=%s", key, types[typeN], url2.QueryEscape(domain))
	val := Common.Tools{}.HttpGet(api)
	fmt.Println("api", api)
	fmt.Println("logs", string(val), types[typeN], key, domain)
	var valResponse Response
	err = json.Unmarshal(val, &valResponse)
	if err != nil {
		return "", err
	}

	if valResponse.Short == "" {
		return "", fmt.Errorf("域名添加失败")
	}

	return valResponse.Short, nil
}
