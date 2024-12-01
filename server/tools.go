package main

import (
	"server/Base"
	"server/Task"
	"time"
)

func main() {

	//status := checkDomain("http://wiudhuhcuhascaocasjhkudhuahfujdf.whvubuu.cn")
	//fmt.Println(status)
	//return
	// 启动初始化
	Base.Base{}.Init()

	go func() {
		for true {
			Task.Day{}.Run()
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for true {
			Task.TimeOut{}.Run()
			time.Sleep(time.Second * 10)
		}
	}()

	for true {
		Task.CheckDomain{}.Run()
		time.Sleep(time.Second)
	}

	print("done")
}
