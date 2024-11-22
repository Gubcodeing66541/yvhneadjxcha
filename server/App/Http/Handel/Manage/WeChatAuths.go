package Manage

import "github.com/gin-gonic/gin"

//微信管理
type WeChatAuth struct {

}

//微信公众号列表
func(WeChatAuth) List(c *gin.Context){

}

//启用关闭
func(WeChatAuth) OpenOrClose(c *gin.Context){

}

//新增配置
func(WeChatAuth) Create(c *gin.Context){

}

//修改配置
func(WeChatAuth) Update(c *gin.Context){

}

//删除配置
func(WeChatAuth) Delete(c *gin.Context){

}