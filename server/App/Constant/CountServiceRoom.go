package Constant

type CountServiceRoom struct {
	ServiceId  int    `json:"service_id"`
	AllUser    int    `json:"all_user"`
	AddUser    int    `json:"add_user"`
	OnlineUser int    `json:"online_user"`
	CountTime  string `json:"count_time"`
}
