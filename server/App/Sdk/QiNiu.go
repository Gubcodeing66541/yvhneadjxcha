package Sdk

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"server/Base"
)
type QiNiu struct {}


func (QiNiu) GetUpToken() string {
	accessKey := Base.AppConfig.Oss.Qiniu.AccessKey
	secretKey := Base.AppConfig.Oss.Qiniu.SecretKey
	bucket := Base.AppConfig.Oss.Qiniu.Bucket
	putPolicy := storage.PutPolicy{Scope: bucket}
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

