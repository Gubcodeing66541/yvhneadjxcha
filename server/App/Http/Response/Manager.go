package Response

type ManagerPayCont struct {
	RenewServiceManager int `json:"renew_service_manager"`
	AllAccount          int `json:"all_account"`
	AllPay              int `json:"all_pay"`
	TodayPay            int `json:"today_pay"`
}
