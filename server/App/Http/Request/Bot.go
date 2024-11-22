package Request

type BotUpdate struct {
	Status string `json:"status"  uri:"status" xml:"status" form:"status"` // stop停止 run启动
	Head   string `json:"head"  uri:"head" xml:"head" form:"head"`
	Hello  string `json:"hello"  uri:"hello" xml:"hello" form:"hello"`
}
