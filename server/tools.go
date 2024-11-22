package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"server/Base"
	"server/Task"
	"strings"
	"time"
)

func main() {

	//status := checkDomain("http://wiudhuhcuhascaocasjhkudhuahfujdf.whvubuu.cn")
	//fmt.Println(status)
	//return
	// 启动初始化
	Base.Base{}.Init()

	go func() {
		for true {
			Task.Day{}.Run()
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for true {
			Task.TimeOut{}.Run()
			time.Sleep(time.Second * 10)
		}
	}()

	for true {
		Task.CheckDomain{}.Run()
		time.Sleep(time.Second * 2)
	}

	print("done")
}

func checkDomains(domain string) bool {
	key := Base.AppConfig.Check.Key
	checkUrl := fmt.Sprintf("https://wx.03426.com/api.php?token=%s&url=%s&type=1", key, domain)
	request, _ := http.NewRequest("GET", checkUrl, nil)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("无法检测域名", domain)
		return true
	}
	page, _ := ioutil.ReadAll(resp.Body)
	val := string(page)
	fmt.Println(val)
	return !(strings.Index(val, "封禁") >= 0)
}
