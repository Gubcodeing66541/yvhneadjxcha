package Task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Model/Common"
	"server/Base"
	"strings"
	"sync"
	"time"
)

type CheckDomain struct{}

func (c CheckDomain) Run() {
	fmt.Println("")
	var list []Common.Domain
	Base.MysqlConn.Find(&list, "status = ?", "enable")
	fmt.Println("执行域名检测本次任务", time.Now(), list)

	var count = 0
	wg := sync.WaitGroup{}
	for _, val := range list {
		count++
		wg.Add(1)

		//wg.Add(1)
		//if count == 3 {
		count = 0
		time.Sleep(time.Second)
		//}
		go func(vals Common.Domain, wg *sync.WaitGroup) {
			defer wg.Done()
			// 域名检测如果被封禁 下架域名并自动绑定已有的域名
			status := c.checkDomain(vals.Domain)
			if status == false {
				tempServiceId := vals.BindServiceId
				Base.MysqlConn.Model(&vals).Updates(map[string]interface{}{"bind_service_id": 0, "we_chat_ban_status": "1", "status": "un_enable"})
				_ = Logic.Domain{}.Bind(tempServiceId)

				// 推送域名封禁提示
				if vals.BindServiceId != 0 {

				}
				pararm := fmt.Sprintf("?service_id=%d&type=%s&content=%s", tempServiceId, "ban", vals.Domain)
				Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + pararm)
			}
			fmt.Println("check domain:", vals.Domain, " STATUS:", status)
		}(val, &wg)
	}

	time.Sleep(time.Second)

	// 健康修复
	//Base.MysqlConn.Find(&list, "status = ?", "un_enable")
	//fmt.Println("执行域名检测本次任务 修复列表", time.Now(), list)
	//
	//count = 0
	////wg := sync.WaitGroup{}
	//for _, val := range list {
	//	wg.Add(1)
	//	count++
	//
	//	//wg.Add(1)
	//	if count == 30 {
	//		count = 0
	//		time.Sleep(time.Second * 2)
	//	}
	//	go func(vals Common.Domain, wg *sync.WaitGroup) {
	//		defer wg.Done()
	//		// 域名检测如果被封禁 下架域名并自动绑定已有的域名
	//		status := c.checkDomain(vals.Domain)
	//		if status == true {
	//			Base.MysqlConn.Model(&vals).Updates(map[string]interface{}{"bind_service_id": 0, "we_chat_ban_status": "success", "status": "enable"})
	//		}
	//		fmt.Println("check domain:", vals.Domain, " STATUS:", status)
	//	}(val, &wg)
	//}
	wg.Wait()
	time.Sleep(time.Second * 2)
}

func (CheckDomain) checkDomain(domain string) bool {
	key := Base.AppConfig.Check.Key
	checkUrl := fmt.Sprintf("http://wx.03426.com/api.php?token=%s&url=%s&type=1", key, domain)
	request, _ := http.NewRequest("GET", checkUrl, nil)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("无法检测域名", domain)
		return true
	}
	page, _ := ioutil.ReadAll(resp.Body)
	val := string(page)
	fmt.Println(val)
	return !(strings.Index(val, "已被封") >= 0)
}
