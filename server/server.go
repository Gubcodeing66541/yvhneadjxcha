package main

import (
	"bufio"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"server/App"
	"server/App/Http/Logic"
	"server/App/Model/Common"
	"server/Base"
	_ "server/docs" // 这里需要引入本地已生成文档
	"time"
)

// @title 客服系统API文档`
// @version 1.0`
// @description 客服系统api `
// @description 客服后端：`
func main() {

	//启动初始化
	Base.Base{}.Init()

	initSqlDate()

	//trace.Start(os.Stderr)
	//defer trace.Stop()

	// 清除缓存
	Logic.System{}.ClearCache()

	// 启动web服务
	HttpServer := gin.Default()

	if Base.AppConfig.Debug {
		pprof.Register(HttpServer)
	}

	App.HttpRoute{}.BindRoute(HttpServer)

	if Base.AppConfig.Debug {
		HttpServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	//启动服务
	_ = HttpServer.Run(":80")
}

func initSqlDate() {
	var rename Common.Rename
	err := Base.MysqlConn.First(&rename).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("initSqlDate", err.Error())
		return
	}
	if rename.Id != 0 {
		fmt.Println("RENAME", rename)
		return
	}

	filename := "./rename.md"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("open file err", err.Error())
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		lineText := scanner.Text() // 获取当前行的文本内容
		fmt.Println(lineText)      // 输出每一行的内容
		renameData := Common.Rename{
			Rename: lineText,
		}
		Base.MysqlConn.Create(&renameData)
		if count >= 500 {
			count = 0
			time.Sleep(time.Second)
		}
		count++
		fmt.Println(lineText)
	}

	if err := scanner.Err(); err != nil {
		panic("读取命名文件时发生错误")
	}
}
