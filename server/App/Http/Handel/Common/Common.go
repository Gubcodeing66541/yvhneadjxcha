package Common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	Common3 "server/App/Common"
	"server/App/Http/Logic"
	"server/App/Http/Request"
	"server/App/Http/Response"
	Common2 "server/App/Model/Common"
	"server/App/Model/Service"
	"server/Base"
	"time"
)

type Common struct{}

func (Common) WeChatFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" || filename == "check" {
		c.String(http.StatusOK, "")
	}
	var res Common2.WeChatAuth
	Base.MysqlConn.Model(&Common2.WeChatAuth{}).Find(&res, "file_name = ?", filename)
	c.String(http.StatusOK, res.FileValue)
}

// @summary 统计
// @tags 客服系统
// @Param token header string true "认证token"
// @Router /service/common/count [post]
func (Common) Count(c *gin.Context) {
	var req Request.ActivateLogin
	err := c.ShouldBind(&req)
	if err != nil {
		Common3.ApiResponse{}.Error(c, "请求繁忙", gin.H{})
		return
	}

	var count []Response.ServiceUserCount

	var service Service.Service
	Base.MysqlConn.Find(&service, "username = ?", req.Username)
	if service.Id == 0 {
		Common3.ApiResponse{}.Success(c, "no", gin.H{"count": count})
		return
	}

	//sql := "select *,ROUND(user_cnt/all_cnt* 100)  as rate from (select  count(*) as all_cnt, count(case  when late_role = 'user' then late_role end) as user_cnt, DATE_FORMAT(create_time,'%Y-%m-%d') as dates from service_rooms where service_id= ? group by  dates order by dates desc) t"

	//var count []Response.ServiceUserCount
	sql := "select count(*) as user_cnt,sum(same_day_is_reply) as msg_cnt,dates,round((sum(same_day_is_reply)/count(*))*100) as rate from (  select t1.user_id,t1.dates,if(t2.user_id is null,0,1) as same_day_is_reply " +
		"from (   select user_id,DATE_FORMAT(create_time,'%Y-%m-%d') as dates from service_rooms where service_id = ? ) t1 " +
		"left join (   select distinct user_id,DATE_FORMAT(create_time,'%Y-%m-%d') as dates  " +
		"from messages where send_role = 'service' and service_id = ?  ) t2 on t1.user_id = t2.user_id and t1.dates = t2.dates) t " +
		"group by  dates order by  dates desc limit 7 "
	Base.MysqlConn.Raw(sql, service.ServiceId, service.ServiceId).Scan(&count)

	//Base.MysqlConn.Raw(sql, service.ServiceId).Scan(&count)
	Common3.ApiResponse{}.Success(c, "ok", gin.H{"count": count})
}

// @summary 获取上传参数
// @tags 客服系统,用户端
// @Param token header string true "认证token"
// @Router /service/common/api/oss_config [post]
func (Common) Oss(c *gin.Context) {
	var req Request.FileName
	err := c.ShouldBind(&req)
	if err != nil {
		Common3.ApiResponse{}.Error(c, "请求参数不完整", gin.H{})
		return
	}
	req.FileName = fmt.Sprintf("image/%d-%d-%d-%s", Common3.Tools{}.GetRoleId(c), time.Now().Unix(), rand.Intn(999999), req.FileName)

	var cnf Common2.SystemConfig
	Base.MysqlConn.Find(&cnf)

	if cnf.OssStorage == "tencent" {
		token, SecretId, SecretKey, err := Logic.Oss{}.GetTencentToken(req.FileName, cnf)
		if err != nil {
			Common3.ApiResponse{}.Error(c, "解析失败", gin.H{"config": cnf, "token": err.Error()})
			return
		}
		Common3.ApiResponse{}.Success(c, "获取成功", gin.H{
			"url":        token,
			"token":      token,
			"secret_id":  SecretId,
			"secret_key": SecretKey,
			"file_name":  fmt.Sprintf("%s/%s", cnf.OssTencentDomain, req.FileName),
		})
		return
	}

	if cnf.OssStorage == "tencent2" {
		token, SecretId, SecretKey, err := Logic.Oss{}.GetTencentToken2(req.FileName, cnf)
		if err != nil {
			Common3.ApiResponse{}.Error(c, "解析失败", gin.H{"config": cnf, "token": err.Error()})
			return
		}
		Common3.ApiResponse{}.Success(c, "获取成功", gin.H{
			"url":        token,
			"token":      token,
			"secret_id":  SecretId,
			"secret_key": SecretKey,
			"file_name":  fmt.Sprintf("%s/%s", cnf.OssTencent2Domain, req.FileName),
		})
		return
	}

	if cnf.OssStorage == "ali" {
		aliToken, err := Logic.Oss{}.GetAliToken(cnf.OssAliRegion, cnf.OssAliAccessKeyId, cnf.OssAliAccessKeySecret)
		if err != nil {
			Common3.ApiResponse{}.Error(c, "解析失败", gin.H{"config": cnf, "token": err.Error()})
			return
		}
		Common3.ApiResponse{}.Success(c, "获取成功", gin.H{"config": cnf, "token": aliToken})
		return
	}

	Common3.ApiResponse{}.Success(c, "获取成功", gin.H{"config": cnf, "token": "system"})
}

// 获取头像图片
func (Common) GetImg(c *gin.Context) {
	var imgUrl string

	fileUrl := "C:/Users/Administrator/Documents/WeChat Files/wxid_7tcicr2yv19u22/FileStorage/File/2022-12/头像"
	files, err := ioutil.ReadDir(fileUrl)
	if err != nil {
		Common3.ApiResponse{}.Error(c, "获取图片失败", gin.H{})
		return
	}

	//初始值
	redisImgCount := Common3.RedisTools{}.GetInt("img_cuts")
	if redisImgCount == -1 || redisImgCount == 0 {
		Common3.RedisTools{}.SetInt("img_cuts", 0+1)
		Common3.ApiResponse{}.Success(c, "获取图片第一个", gin.H{"img": files[0].Name()})
		return
	}

	//是否到头像大小,重0开始继续 不然就继续取下一个头像

	if redisImgCount >= len(files) {
		imgUrl = fmt.Sprintf("/home/head/%s", files[0].Name())
		Common3.RedisTools{}.SetInt("img_cuts", 0)
	} else {
		Common3.RedisTools{}.SetInt("img_cuts", redisImgCount+1)
		imgUrl = fmt.Sprintf("/home/head/%s", files[redisImgCount].Name())
	}

	//获取到最新头像
	Common3.ApiResponse{}.Success(c, "获取图片", gin.H{"img": imgUrl})
}

// 获取随机昵称
func (Common) GetRename(c *gin.Context) {
	var rename string

	//初始值
	renameCount := Common3.RedisTools{}.GetInt("rename_cut")
	if renameCount == -1 {
		Common3.RedisTools{}.SetInt("rename_cut", 1) //初始值
		Common3.ApiResponse{}.Success(c, "获取图片第一个", gin.H{"img": renameCount})
		return
	}

	//查数据库rename昵稱
	var renameDb Common2.Rename
	Base.MysqlConn.Find(&renameDb, "id=?", renameCount)
	if renameDb.Id == 0 {
		var newRename Common2.Rename
		Base.MysqlConn.Find(&newRename, "id=1")
		Common3.RedisTools{}.SetInt("rename_cut", 1)
		rename = newRename.Rename
	} else {
		Common3.RedisTools{}.SetInt("rename_cut", renameCount+1)
		rename = renameDb.Rename
	}

	Common3.ApiResponse{}.Success(c, "获取rename", gin.H{"img": rename})
}

func (Common) CreateRename(c *gin.Context) {

	Common3.ApiResponse{}.Success(c, "获取rename", gin.H{})
}
