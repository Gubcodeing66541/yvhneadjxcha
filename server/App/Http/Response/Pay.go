package Response

type PayCount struct {
	RenewServiceManager int `json:"renew_service_manager"`
	RenewService        int `json:"renew_service"`
	CreateService       int `json:"create_service"`
}
