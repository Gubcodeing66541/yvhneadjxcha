package Constant

type GroupMessageMq struct {
	GroupId   int    `json:"group_id"`
	Type string `json:"type"`
	Info string `json:"info"`
	SendRole string `json:"send_role"`
}