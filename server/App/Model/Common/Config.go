package Common

type SystemConfig struct {
	Id                    int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SystemName            string `json:"system_name"`
	SystemVideo           string `json:"system_video"`
	SystemVersion         string `json:"system_version"`
	SystemLogo            string `json:"system_logo"`
	PcDownloadLink        string `json:"pc_download_link"`
	OssStorage            string `json:"oss_storage"`
	OssAliAccessKeyId     string `json:"oss_ali_access_key_id"`
	OssAliAccessKeySecret string `json:"oss_ali_access_key_secret"`
	OssAliStorage         string `json:"oss_ali_storage"`
	OssAliDomain          string `json:"oss_ali_domain"`
	OssAliRegion          string `json:"oss_ali_region"`
	OssTencentAppid       string `json:"oss_tencent_appid" uri:"oss_tencent_appid" xml:"oss_tencent_appid" form:"oss_tencent_appid"`
	OssTencentSecretId    string `json:"oss_tencent_secret_id" uri:"oss_tencent_secret_id" xml:"oss_tencent_secret_id" form:"oss_tencent_secret_id"`
	OssTencentSecretKey   string `json:"oss_tencent_secret_key" uri:"oss_tencent_secret_key" xml:"oss_tencent_secret_key" form:"oss_tencent_secret_key"`
	OssTencentStorage     string `json:"oss_tencent_storage"  uri:"oss_tencent_storage" xml:"oss_tencent_storage" form:"oss_tencent_storage"`
	OssTencentRegion      string `json:"oss_tencent_region" uri:"oss_tencent_region" xml:"oss_tencent_region" form:"oss_tencent_region"`
	OssTencentDomain      string `json:"oss_tencent_domain" uri:"oss_tencent_domain" xml:"oss_tencent_domain" form:"oss_tencent_domain"`
	OssTencent2Appid      string `json:"oss_tencent2_appid" uri:"oss_tencent2_appid" xml:"oss_tencent2_appid" form:"oss_tencent2_appid"`
	OssTencent2SecretId   string `json:"oss_tencent2_secret_id" uri:"oss_tencent2_secret_id" xml:"oss_tencent2_secret_id" form:"oss_tencent2_secret_id"`
	OssTencent2SecretKey  string `json:"oss_tencent2_secret_key" uri:"oss_tencent2_secret_key" xml:"oss_tencent2_secret_key" form:"oss_tencent2_secret_key"`
	OssTencent2Storage    string `json:"oss_tencent2_storage"  uri:"oss_tencent2_storage" xml:"oss_tencent2_storage" form:"oss_tencent2_storage"`
	OssTencent2Region     string `json:"oss_tencent2_region" uri:"oss_tencent2_region" xml:"oss_tencent2_region" form:"oss_tencent2_region"`
	OssTencent2Domain     string `json:"oss_tencent2_domain" uri:"oss_tencent2_domain" xml:"oss_tencent2_domain" form:"oss_tencent2_domain"`
	AdDefault             string `json:"ad_default"` // 隐私
	PayWechatMchId        string `json:"pay_wechat_mch_id"`
	PayWechatAppid        string `json:"pay_wechat_appid"`
	PayWechatAppKey       string `json:"pay_wechat_app_key"`
	PayAliMchId           string `json:"pay_ali_mch_id"`
	PayAliAppid           string `json:"pay_ali_appid"`
	PayAliAppKey          string `json:"pay_ali_app_key"`
	Pay                   int    `json:"pay"`
	RealName              int    `json:"real_name"`
}
