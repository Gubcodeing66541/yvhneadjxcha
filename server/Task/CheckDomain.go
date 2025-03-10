package Task

import (
	"fmt"
	"io/ioutil"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strings"
	"time"
)

type CheckDomain struct{}

func (c CheckDomain) Run() {
	fmt.Println("")

	// 获取入口和落地
	action := Logic.Domain{}.GetAction()
	transfer := Logic.Domain{}.GetTransfer()
	publicDomain5Number := Logic.Domain{}.GetPublicBindDomain()
	domain := []string{action, transfer.Domain}

	for _, val := range publicDomain5Number {
		time.Sleep(time.Second)
		status := c.checkDomain(val)
		if status == false {
			Base.MysqlConn.Model(&Common.Domain{}).
				Where("domain = ?", val).Update("status", "un_enable")

			// 给所有在线的service发送消息
			pararm := fmt.Sprintf("?service_id=%d&type=%s&content=%s", 0, "ban", val)
			Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + pararm)
		}
	}

	fmt.Println("执行域名检测本次任务", time.Now(), domain)
	for _, val := range domain {
		time.Sleep(time.Second)
		status := c.checkDomain(val)
		if status == false {
			//pararm := fmt.Sprintf("?service_id=%d&type=%s&content=%s", tempServiceId, "ban", val.Domain)
			//Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + pararm)
			Base.MysqlConn.Model(&Common.Domain{}).
				Where("domain = ?", val).Update("status", "un_enable")
		}
		fmt.Println("check domain: join and transfer", val, " STATUS:", status)
	}

	time.Sleep(time.Second)
	var domains []Service2.Service
	Base.MysqlConn.Find(&domains, "time_out > ?", time.Now())

	for _, valInfo := range domains {
		time.Sleep(time.Second)
		status := c.checkDomain(valInfo.Domain)
		if status == false {
			domainInfo := Logic.Domain{}.GetPublic()
			u := fmt.Sprintf("%s?code=%s", domainInfo, valInfo.Code)
			//domainInfo := Logic.Domain{}.GetTransfer()
			//web := fmt.Sprintf("%s/user/auth/local_storage/join_new?code=%s", domainInfo.Domain, valInfo.Code)
			//u, err := Sdk.CreateDomain(Base.AppConfig.DomainKey, web)
			//if err != nil {
			//	fmt.Println("域名创建失败", valInfo.Domain, " STATUS:", status)
			//} else {
			//	Base.MysqlConn.Model(&Service2.Service{}).Where("service_id =?", valInfo.ServiceId).
			//		Update("domain", u)
			//}

			Base.MysqlConn.Model(&Service2.Service{}).Where("service_id =?", valInfo.ServiceId).
				Update("domain", u)

			pararm := fmt.Sprintf("?service_id=%d&type=%s&content=%s", valInfo.ServiceId, "ban", valInfo.Domain)
			Common2.Tools{}.HttpGet("http://127.0.0.1/api/socket/send_to_service_socket" + pararm)

			// 更换域名
		}
		fmt.Println("check domain: join and transfer", valInfo.Domain, " STATUS:", status)
	}

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
	//wg.Wait()
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
	// if !(strings.Index(val, "未知错误") >= 0) {
	// 	return false
	// }
	return !(strings.Index(val, "已被封") >= 0)
}
