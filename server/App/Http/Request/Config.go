package Request

type SystemConfig struct {
	SystemName            string `json:"system_name" uri:"system_name" xml:"system_name" form:"system_name"`
	SystemVideo           string `json:"system_video" uri:"system_video" xml:"system_video" form:"system_video"`
	SystemLogo            string `json:"system_logo" uri:"system_logo" xml:"system_logo" form:"system_logo"`
	SystemVersion         string `json:"system_version" uri:"system_version" xml:"system_version" form:"system_version"`
	PcDownloadLink        string `json:"pc_download_link" uri:"pc_download_link" xml:"pc_download_link" form:"pc_download_link"`
	OssStorage            string `json:"oss_storage"  uri:"oss_storage" xml:"oss_storage" form:"oss_storage"`
	OssAliAccessKeyId     string `json:"oss_ali_access_key_id" uri:"oss_ali_access_key_id" form:"oss_ali_access_key_id" xml:"oss_ali_access_key_id"`
	OssAliAccessKeySecret string `json:"oss_ali_access_key_secret" uri:"oss_ali_access_key_secret" form:"oss_ali_access_key_secret" xml:"oss_ali_access_key_secret"`
	OssAliStorage         string `json:"oss_ali_storage" uri:"oss_ali_storage" form:"oss_ali_storage" xml:"oss_ali_storage"`
	OssAliDomain          string `json:"oss_ali_domain" uri:"oss_ali_domain" form:"oss_ali_domain" xml:"oss_ali_domain"`
	OssAliRegion          string `json:"oss_ali_region" uri:"oss_ali_region" form:"oss_ali_region" xml:"oss_ali_region"`
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
	AdDefault             string `json:"ad_default" uri:"ad_default" xml:"ad_default" form:"ad_default"`
	PayWechatMchId        string `json:"pay_wechat_mch_id" uri:"pay_wechat_mch_id" xml:"pay_wechat_mch_id" form:"pay_wechat_mch_id"`
	PayWechatAppid        string `json:"pay_wechat_appid" uri:"pay_wechat_appid" xml:"pay_wechat_appid" form:"pay_wechat_appid"`
	PayWechatAppKey       string `json:"pay_wechat_app_key" uri:"pay_wechat_app_key" xml:"pay_wechat_app_key" form:"pay_wechat_app_key"`
	PayAliMchId           string `json:"pay_ali_mch_id" uri:"pay_ali_mch_id" xml:"pay_ali_mch_id" form:"pay_ali_mch_id"`
	PayAliAppid           string `json:"pay_ali_appid" uri:"pay_ali_appid" xml:"pay_ali_appid" form:"pay_ali_appid"`
	PayAliAppKey          string `json:"pay_ali_app_key" uri:"pay_ali_app_key" xml:"pay_ali_app_key" form:"pay_ali_app_key"`
	Pay                   int    `json:"pay" uri:"pay" xml:"pay" form:"pay"`
	RealName              int    `json:"real_name" uri:"real_name" xml:"real_name" form:"real_name"`
}
