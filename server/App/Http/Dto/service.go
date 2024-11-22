package Dto

type ServiceLoginInfo struct {
	ServiceId int    `json:"service_id"`
	Member    string `json:"member"`
	Password  string `json:"password"`
}
