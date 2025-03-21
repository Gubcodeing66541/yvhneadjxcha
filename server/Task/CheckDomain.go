package Task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	Common2 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
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

		time.Sleep(time.Second)
		status = c.checkDomain(action)
		if status == false {
			action = Logic.Domain{}.GetAction()
			Base.MysqlConn.Model(&Common.Domain{}).
				Where("domain = ?", action).Update("status", "un_enable")
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
	// if !(strings.Index(val, "未知错误") >= 0) {
	// 	return false
	// }

	type resa struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var res resa
	_ = json.Unmarshal([]byte(val), &res)

	fmt.Println(domain, val, res)

	if res.Code == 200 {
		return true
	}
	return false

	//return !(strings.Index(val, "已被封") >= 0)
}
