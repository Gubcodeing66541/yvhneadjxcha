package Service

type ServiceMenuSetting struct {
	Id      int    `form:"id" json:"id" uri:"id" xml:"id" `
	Title   string `form:"title" json:"title" uri:"title" xml:"title" `
	Content string `form:"content" json:"content" uri:"content" xml:"content" `
	Action  string `form:"action" json:"action" uri:"action" xml:"action" `
	Sort    int    `form:"sort" json:"sort" uri:"sort" xml:"sort" `
	Tag     string `form:"tag" json:"tag" uri:"tag" xml:"tag" `
}
