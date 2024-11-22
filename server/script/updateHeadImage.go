package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

func main() {
	// 注册用户
	path := "C:\\Users\\Administrator\\Documents\\WeChat Files\\wxid_7tcicr2yv19u22\\FileStorage\\File\\2022-12\\头像"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("open error")
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())

	fmt.Println(len(files))
	fmt.Printf(files[0].Name())
	fmt.Println()
	fmt.Printf(files[len(files)-1].Name())
	fmt.Println("开始rename")
	for i := 0; i < len(files)-1; i++ {
		pathTempm := path + "\\" + files[i].Name()
		newfileName := fmt.Sprintf("head%d.jpg", i)
		err := os.Rename(pathTempm, fmt.Sprintf("head%d.jpg", i))
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("p:", pathTempm, "n", newfileName)
	}
	fmt.Println("结束rename")
}
