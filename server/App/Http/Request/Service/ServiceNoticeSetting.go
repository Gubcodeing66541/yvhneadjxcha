package Service

type ServiceNoticeSetting struct {
	IsShow int    `json:"is_show" form:"is_show" json:"is_show" uri:"is_show" xml:"is_show" `
	Image  string `json:"image" form:"image" json:"image" uri:"image" xml:"image" `
	Text   string `json:"text" form:"text" json:"text" uri:"text" xml:"text" `
}
