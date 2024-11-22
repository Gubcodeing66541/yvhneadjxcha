package Logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/tencentyun/cos-go-sdk-v5"
	"hash"
	"io"
	"net/http"
	"net/url"
	Common2 "server/App/Model/Common"
	"time"
)

type Oss struct{}

// 获取腾讯云oss
func (Oss) GetTencentToken(fileName string, cnf Common2.SystemConfig) (token string, SecretId string, SecretKey string, err error) {

	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		cnf.OssTencentStorage, cnf.OssTencentRegion))

	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cnf.OssTencentSecretId,
			SecretKey: cnf.OssTencentSecretKey,
		},
	})

	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPut, fileName, cnf.OssTencentSecretId, cnf.OssTencentSecretKey, time.Hour, nil)
	if err != nil {
		panic(err)
	}

	return presignedURL.String(), SecretId, SecretKey, nil
}

// 获取腾讯云oss
func (Oss) GetTencentToken2(fileName string, cnf Common2.SystemConfig) (token string, SecretId string, SecretKey string, err error) {

	// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com",
		cnf.OssTencent2Storage, cnf.OssTencent2Region))

	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cnf.OssTencent2SecretId,
			SecretKey: cnf.OssTencent2SecretKey,
		},
	})

	ctx := context.Background()

	// 获取预签名URL
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPut, fileName, cnf.OssTencent2SecretId, cnf.OssTencent2SecretKey, time.Hour, nil)
	if err != nil {
		panic(err)
	}

	return presignedURL.String(), SecretId, SecretKey, nil
}

// 获取阿里云oss
func (Oss) GetAliToken(OssAliRegion string, OssAliAccessKeyId string, OssAliAccessKeySecret string) (string, error) {
	//构建一个阿里云客户端, 用于发起请求。
	//设置调用者（RAM用户或RAM角色）的AccessKey ID和AccessKey Secret。
	//第一个参数就是bucket所在位置，可查看oss对象储存控制台的概况获取
	//第二个参数就是步骤一获取的AccessKey ID
	//第三个参数就是步骤一获取的AccessKey Secret
	client, err := sts.NewClientWithAccessKey(OssAliRegion, OssAliAccessKeyId, OssAliAccessKeySecret)

	//构建请求对象。
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	//设置参数。关于参数含义和设置方法，请参见《API参考》。
	//request.RoleArn = "<RoleArn>"                 //步骤三获取的角色ARN
	//request.RoleSessionName = "<RoleSessionName>" //步骤三中的RAM角色名称

	//发起请求，并得到响应。
	response, err := client.AssumeRole(request)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 结构体
	type Policy struct {
		Expiration string          `json:"expiration"`
		Conditions [][]interface{} `json:"conditions"`
	}
	// 生成签名代码
	var policy Policy
	policy.Expiration = "9999-12-31T12:00:00.000Z"
	var conditions []interface{}
	conditions = append(conditions, "content-length-range")
	conditions = append(conditions, 0)
	conditions = append(conditions, 1048576000)
	policy.Conditions = append(policy.Conditions, conditions)
	policyByte, err := json.Marshal(policy)
	if err != nil {
		return "", errors.New("序列化失败")
	}
	policyBase64 := base64.StdEncoding.EncodeToString(policyByte)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(response.Credentials.AccessKeySecret))
	io.WriteString(h, policyBase64)
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	// 主要拿的就是下面这两个玩意
	fmt.Println("policyBase64：", policyBase64)
	fmt.Println("signature", signature)
	// 将这两个与前面获取的临时授权参数一起返回就好了
	return signature, errors.New("序列化失败")
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code": 1,
	//	"data": gin.H{
	//		"response":     response,
	//		"policyBase64": policyBase64,
	//		"signature":    signature,
	//	},
	//	"msg": "获取成功",
	//})

}
