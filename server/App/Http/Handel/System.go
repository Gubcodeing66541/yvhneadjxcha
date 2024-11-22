package Handel

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"server/App/Common"
	"server/App/Http/Logic"
	"server/Base"
	"time"
)

type System struct{}

func (System) Status(c *gin.Context) {
	Common.ApiResponse{}.Success(c, "ok", gin.H{
		"online_count": Base.WebsocketHub.GetOnlineCount(),
		"system_info":  Base.WebsocketHub.GetAllStatus(),
	})
}

func (s System) ClearCache(c *gin.Context) {
	Logic.System{}.ClearCache()
	Common.ApiResponse{}.Success(c, "OK", gin.H{})
}

func (s System) UploadImage(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		Common.ApiResponse{}.Error(c, "请选择需要上传的文件", gin.H{"err": err.Error()})
	}

	FileMap := map[string]string{"image/png": "png", "image/gif": "gif", "image/jpeg": "jpg", "video/mp4": "mp4"}
	if FileMap[f.Header.Get("Content-Type")] == "" {
		Common.ApiResponse{}.Error(c, "只允许上传png gif jpg格式图片和mp4格式的视频", gin.H{})
		return
	}

	if f.Size > 1024*1024*60 {
		Common.ApiResponse{}.Error(c, "文件大小不能超出60mb", gin.H{"size": f.Size})
		return
	}

	// 拼接保存路径 将读取的文件保存在服务端
	rootPath := fmt.Sprintf("/static/upload/image/%s", time.Now().Format("20060102"))
	fileName := fmt.Sprintf("%s.%s", Common.Tools{}.RandFileName(c), FileMap[f.Header.Get("Content-Type")])
	dst := path.Join("."+rootPath, fileName)
	_ = os.MkdirAll("."+rootPath, os.ModePerm)

	// 保存文件
	err = c.SaveUploadedFile(f, dst)
	if err != nil {
		Common.ApiResponse{}.Error(c, "error", gin.H{"err": err.Error(), "file": dst})
		return
	}

	filePath := Base.AppConfig.HttpHost + rootPath + "/" + fileName
	Common.ApiResponse{}.Success(c, "OK", gin.H{"file_name": filePath, "file_type": FileMap[f.Header.Get("Content-Type")]})
}

// @summary 系统默认文件上传
// @tags 公共接口
// @Param token header string true "认证token"
// @Param image  formData file true "文件参数"
// @Router /api/system/upload [post]
func (s System) Upload(c *gin.Context) {
	f, err := c.FormFile("image")
	if err != nil {
		Common.ApiResponse{}.Error(c, "请选择需要上传的文件", gin.H{"err": err.Error()})
	}

	FileMap := map[string]string{"image/png": "png", "image/gif": "gif", "image/jpeg": "jpg", "video/mp4": "mp4"}
	if FileMap[f.Header.Get("Content-Type")] == "" {
		Common.ApiResponse{}.Error(c, "只允许上传png gif jpg格式图片和mp4格式的视频", gin.H{})
		return
	}

	if f.Size > 1024*1024*60 {
		Common.ApiResponse{}.Error(c, "文件大小不能超出60mb", gin.H{"size": f.Size})
		return
	}

	// 拼接保存路径 将读取的文件保存在服务端
	rootPath := fmt.Sprintf("/static/upload/%s", time.Now().Format("20060102"))
	fileName := fmt.Sprintf("%s.%s", Common.Tools{}.RandFileName(c), FileMap[f.Header.Get("Content-Type")])
	dst := path.Join("."+rootPath, fileName)
	_ = os.MkdirAll("."+rootPath, os.ModePerm)

	// 保存文件
	err = c.SaveUploadedFile(f, dst)
	if err != nil {
		Common.ApiResponse{}.Error(c, "error", gin.H{"err": err.Error(), "file": dst})
		return
	}

	filePath := Base.AppConfig.HttpHost + rootPath + "/" + fileName
	Common.ApiResponse{}.Success(c, "OK", gin.H{"file_name": filePath, "file_type": FileMap[f.Header.Get("Content-Type")]})
}

func (s System) Action(c *gin.Context) {
	c.String(http.StatusOK, Logic.Domain{}.GetAction())
	//Common.ApiResponse{}.Success(c, "未知操作", gin.H{"action": Logic.Domain{}.GetAction()})
}
