package Request

type SendToUser struct {
	To   int    `form:"to" json:"to" uri:"to" xml:"to" binding:"required"`
	Type string `form:"type" json:"type" uri:"type" xml:"type" binding:"required"`
	Info string `form:"info" json:"info" uri:"info" xml:"info" binding:"required"`
}

type MsgPage struct {
	Id int `form:"id" json:"id" uri:"id" xml:"id"`
}

type MsgUserId struct {
	UserId int `json:"user_id"`
}

type MsgList struct {
	ServiceId int `json:"service_id" uri:"service_id" xml:"service_id" form:"service_id"`
	UserId    int `json:"user_id" uri:"user_id" xml:"user_id" form:"user_id"`
	Page      int `json:"page" uri:"page" xml:"page" form:"page"`
	Offset    int `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type ServiceSendMessage struct {
	UserId  int    `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Type    string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
}

type ServiceSendMessageGroup struct {
	UserId  []int  `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Type    string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
}

type UserSendMessage struct {
	Type           string `json:"type" form:"type" json:"type" uri:"type" xml:"type"`
	Content        string `json:"content" form:"content" json:"content" uri:"content" xml:"content"`
	ServiceContent string `json:"service_content" form:"service_content" uri:"service_content" xml:"service_content"`
}

type UpdateServiceDetail struct {
	Head           string `json:"head" form:"head" json:"head" uri:"head" xml:"head"`
	Name           string `json:"name" form:"name" json:"name" uri:"name" xml:"name"`
	CodeBackground string `json:"code_background" form:"code_background" json:"code_background" uri:"code_background" xml:"code_background"`
	CodeColor      string `json:"code_color" form:"code_color" json:"code_color" uri:"code_color" xml:"code_color"`
}

type RemoveMsg struct {
	UserId int `json:"user_id" form:"user_id" json:"user_id" uri:"user_id" xml:"user_id"`
	Id     int `json:"id" form:"id" json:"id" uri:"id" xml:"id"`
}
