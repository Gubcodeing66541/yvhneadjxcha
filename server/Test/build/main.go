package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

var headers = map[string]interface{}{
	"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 18_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.54(0x18003638) NetType/WIFI Language/zh_CN",
	"Accept-Encoding": "gzip,compress,br,deflate",
	"token":           "hzfjTokeneyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJoemZqIiwiaWQiOiIxODczMDMzNjkyODM1MjIxNTA1IiwibmFtZSI6IuiSi-WNjuS4nCIsImFjY291bnROYW1lIjoiMTMwNjMzMzMwNzAiLCJhcmVhQ29kZSI6IjM3MDEwMjAwMSIsImlhdCI6MTczNTQwMTAwOSwiZXhwIjoxNzM1NDg3NDA5fQ.fLEMx8etCx7ByXZeN4rfUvr8fI-mpv2k5vwxPK9HWtI",
	"Referer":         "https://servicewechat.com/wx5e6e47eae1827134/24/page-frame.html",
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
			flist = append(flist, GetDomain("join.html", headers))
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
		// Example headers

		for i := 0; i < numI; i++ {
			number++
			if number >= 5 {
				time.Sleep(time.Second)
				number = 0
			}
			flist = append(flist, GetDomain("action.html", headers))
		}

		stv := ""
		for _, v := range flist {
			stv += v + "\n"
		}
		c.String(http.StatusOK, stv)

	})
	s.Run(":80")

}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		fmt.Println("Error generating random string:", err)
		return ""
	}
	for i := 0; i < len(result); i++ {
		result[i] = charset[result[i]%byte(len(charset))]
	}
	return string(result)
}

// CreatePost function to upload a file via POST request
func GetDomain(file string, head map[string]interface{}) string {
	// Generate a random filename with .jpeg extension
	randomFileName := GenerateRandomString(10) + ".jpeg"
	fmt.Println("Generated Random Filename:", randomFileName)

	// Create a buffer to hold the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Open the file to be uploaded
	fileHandle, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer fileHandle.Close()

	// Create the form-data file part, using the random filename
	part, err := writer.CreateFormFile("file", randomFileName)
	if err != nil {
		return ""
	}

	// Copy the file content into the form part
	_, err = io.Copy(part, fileHandle)
	if err != nil {
		return ""
	}

	// Add additional fields if needed (e.g., 'fileCode', etc.)
	for key, value := range head {
		_ = writer.WriteField(key, fmt.Sprintf("%v", value))
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return ""
	}

	// Prepare the request
	req, err := http.NewRequest("POST", "https://sdhh.yjt.shandong.gov.cn:9443/weld-standard/hzfj-system/system/file/upload", body)
	if err != nil {
		return ""
	}

	// Set the headers from the input map
	for key, value := range head {
		req.Header.Set(key, fmt.Sprintf("%v", value))
	}

	// Set the Content-Type header to the multipart writer's value
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using the http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Print the response status and body for debugging
	fmt.Println("Response Status:", resp.Status)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	fmt.Println("Response Body:", string(responseBody))

	type responseData struct {
		Data struct {
			Url string `json:"url`
		} `json:"data"`
	}

	var resPon responseData

	json.Unmarshal(responseBody, &resPon)

	return resPon.Data.Url
}

func ResetDomain(fileName string) string {
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
