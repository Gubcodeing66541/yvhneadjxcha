package Logic

import (
	"errors"
	"server/App/Model/Common"
	Service2 "server/App/Model/Service"
	"server/Base"
	"strings"
	"time"
)

type Domain struct{}

func (Domain) GetPublic() string {
	var domain Common.Domain
	Base.MysqlConn.Limit(1).Order("id asc").Find(&domain, "type = ? and status = 'enable'", "private")
	if domain.Domain != "" {
		return domain.Domain
	}
	return ""
}

func (Domain) GetPublicDomain5Number() []string {
	var domain []Common.Domain
	Base.MysqlConn.Limit(5).Order("id asc").Find(&domain, "type = ? and status = 'enable'", "private")
	if len(domain) > 0 {
		domainList := []string{}
		for _, val := range domain {
			domainList = append(domainList, val.Domain)
		}
		return domainList
	}
	return []string{}
}

func (Domain) GetPublicBindDomain() []string {
	var domain []Common.Domain
	Base.MysqlConn.Order("id asc").Find(&domain, "type = ? and status = 'enable'", "private")
	if len(domain) > 0 {
		domainList := []string{}
		for _, val := range domain {
			domainList = append(domainList, val.Domain)
		}
		return domainList
	}

	return []string{}
}

func (Domain) GetAction() string {
	var domain Common.Domain
	Base.MysqlConn.Limit(1).Find(&domain, "`type` = ? and status = 'enable'", "action")
	if domain.Domain != "" {
		return domain.Domain
	}
	return ""
}

func (Domain) Get(id int) Common.Domain {
	var Model Common.Domain
	Base.MysqlConn.Find(&Model, "id = ?", id)
	return Model
}

func (Domain) GetNoUsePrivateNum() int {
	var leng int
	Base.MysqlConn.Model(&Common.Domain{}).Where("`type` = 'private' and bind_service_id = 0 and we_chat_ban_status = 'success' and status = 'enable'").Count(&leng)
	return leng
}

func (Domain) GetTransfer() Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Model(&Common.Domain{}).Order("id asc").Find(&domain, "`type` = 'transfer' and status = 'enable' and we_chat_ban_status = 'success' ")
	return domain
}

func (Domain) GetServiceBind(serviceId int) Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "bind_service_id = ?", serviceId)
	return domain
}

func (Domain) Bind(serviceId int) error {
	err := ServiceIsOverdue(serviceId)
	if err != nil {
		return err
	}

	var service Service2.Service
	Base.MysqlConn.Find(&service, "service_id = ?", serviceId)

	// 检测1，2，3域名是否不能用用了
	domainInfo := Common.Domain{}
	Base.MysqlConn.Find(&domainInfo, "id = ? and `status` = ?", service.BindDomainId, "enable")
	bindDomains := []int{service.BindDomainId, service.BindDomainId2, service.BindDomainId3}
	if domainInfo.Id == 0 {
		Base.MysqlConn.Model(&Service2.Service{}).Where("id = ?", domainInfo.Id).Update("bind_domain_id", 0)
		domainItem := Common.Domain{}
		Base.MysqlConn.Find(&domainItem, "`status` = ? and  `type` = ?  and  bind_cnt < 5 and id not in (?)", "enable", "private", bindDomains)
		if domainItem.Id == 0 {
			return errors.New("无可用分配域名")
		}
		Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", service.Id).Updates(map[string]interface{}{"bind_domain_id": domainItem.Id})
		Base.MysqlConn.Model(&Common.Domain{}).Where("id = ?", domainItem.Id).Update("bind_cnt", domainItem.BindCnt+1)
		Base.MysqlConn.Model(&domainItem).Where("id = ?", domainItem.Id).Update("bind_service_id", serviceId)
		bindDomains = append(bindDomains, domainItem.Id)
	}

	// 检测2域名是否能用
	domainInfo2 := Common.Domain{}
	Base.MysqlConn.Find(&domainInfo2, "id = ? and `status` = ?", service.BindDomainId2, "enable")
	if domainInfo2.Id == 0 {
		Base.MysqlConn.Model(&Service2.Service{}).Where("id = ?", domainInfo2.Id).Update("bind_domain_id2", 0)
		domainItem := Common.Domain{}
		Base.MysqlConn.Find(&domainItem, "`status` = ? and  `type` = ?  and  bind_cnt < 5 and id not in (?)", "enable", "private", bindDomains)
		if domainItem.Id == 0 {
			return errors.New("无可用分配域名")
		}
		Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", service.Id).Updates(map[string]interface{}{"bind_domain_id2": domainItem.Id})
		Base.MysqlConn.Model(&Common.Domain{}).Where("id = ?", domainItem.Id).Update("bind_cnt", domainItem.BindCnt+1)
		Base.MysqlConn.Model(&domainItem).Where("id = ?", domainItem.Id).Update("bind_service_id", serviceId)
		bindDomains = append(bindDomains, domainItem.Id)
	}

	// 检测3域名是否能用
	domainInfo3 := Common.Domain{}
	Base.MysqlConn.Find(&domainInfo3, "id = ? and `status` = ?", service.BindDomainId3, "enable")
	if domainInfo3.Id == 0 {
		Base.MysqlConn.Model(&Service2.Service{}).Where("id = ?", domainInfo3.Id).Update("bind_domain_id3", 0)
		domainItem := Common.Domain{}
		Base.MysqlConn.Find(&domainItem, "`status` = ? and  `type` = ?  and  bind_cnt < 5 and id not in (?)", "enable", "private", bindDomains)
		if domainItem.Id == 0 {
			return errors.New("无可用分配域名")
		}
		Base.MysqlConn.Model(&Service2.Service{}).Where("service_id = ?", service.Id).Updates(map[string]interface{}{"bind_domain_id3": domainItem.Id})
		Base.MysqlConn.Model(&Common.Domain{}).Where("id = ?", domainItem.Id).Update("bind_cnt", domainItem.BindCnt+1)
		Base.MysqlConn.Model(&domainItem).Where("id = ?", domainItem.Id).Update("bind_service_id", serviceId)
	}

	return nil
}

func (Domain) UnEnable(domainId int) {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", domainId)

	if domain.Id == 0 {
		return
	}

	domain.Status = "un_enable"
	domain.BindServiceId = 0
	Base.MysqlConn.Update(&domain)

	//推送消息

}

func (Domain) QueryById(id int) Common.Domain {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)
	return domain
}

func (Domain) Update(id int, domain string, typeEd string, status string) error {
	var domainEntity Common.Domain
	Base.MysqlConn.Find(&domainEntity, "id = ?", id)
	if domainEntity.BindServiceId != 0 {
		return errors.New("该域名已绑定客服，请解绑后再修改")
	}

	Base.MysqlConn.Model(&domainEntity).Where("id = ?", id).Updates(Common.Domain{Domain: domain, Type: typeEd, Status: status})
	return nil
}

func (Domain) Delete(id int) error {
	if id == 0 {
		return errors.New("参数不全")
	}

	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)
	//if domain.BindServiceId != 0 {
	//	return errors.New("该域名已绑定客服，请解绑后再删除")
	//}

	Base.MysqlConn.Delete(&domain, "id = ?", id)
	return nil
}

func (Domain) Create(domain string, typeEd string, status string) {

	// 将字符串 分割成 字符串数组
	// 参数：要拼接的字符串，分割的内容
	domainList := strings.Split(domain, "\n")
	for _, domain = range domainList {

		if !strings.HasPrefix(domain, "http") {
			domain = "http://" + domain
		}
		var domaineEntity Common.Domain
		domaineEntity.Domain = domain
		domaineEntity.CreateTime = time.Now()
		domaineEntity.UpdateTime = time.Now()
		domaineEntity.Status = status
		domaineEntity.WeChatBanStatus = "success"
		if typeEd != "" {
			domaineEntity.Type = typeEd
		} else {
			domaineEntity.Type = "private"
		}
		Base.MysqlConn.Create(&domaineEntity)
	}
}

func ServiceIsOverdue(serviceId int) error {
	var seviceAuth Service2.ServiceAuth
	Base.MysqlConn.Find(&seviceAuth, "service_id = ?", serviceId)
	if seviceAuth.ServiceId == 0 {
		return errors.New("客服不存在")
	}

	now := time.Now()
	if seviceAuth.TimeOut.Before(now) {
		return errors.New("客服过期")
	}
	return nil
}

func (Domain) EnableDisable(id int, status string) {
	var domain Common.Domain
	Base.MysqlConn.Find(&domain, "id = ?", id)

	if domain.BindServiceId > 0 && status == "un_enable" {
		bindServiceId := domain.BindServiceId

		Base.MysqlConn.Model(&domain).Where("id = ?", domain.Id).Updates(&Common.Domain{Status: "un_enable", BindServiceId: 0})

		var domainStruct Domain
		domainStruct.Bind(bindServiceId)
	} else {
		Base.MysqlConn.Model(&domain).Where("id = ?", id).Update(&Common.Domain{Status: "un_enable", BindServiceId: 0})
	}
}

func (d Domain) CheckBindDomain(service Service2.Service) {
	// 检测账号是否过期
	if service.TimeOut.Unix()-time.Now().Unix() <= 0 {
		return
	}

	// 检测域名是否绑定
	Domain{}.Bind(service.ServiceId)
}
