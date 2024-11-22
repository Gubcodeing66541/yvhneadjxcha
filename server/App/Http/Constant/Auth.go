package Constant

import "time"

type UserAuthToken struct {
	RoleId    int
	RoleType  string // service user manage
	RandStr   string
	ServiceId int
	GroupId   int
	Time      time.Time
	Key       string
}
