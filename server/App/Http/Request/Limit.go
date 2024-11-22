package Request

type PageLimit struct {
	Search string `json:"search" uri:"search" xml:"search" form:"search"`
	Page   int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}

type PageLimitByDate struct {
	StartTime string `json:"start_time" uri:"start_time" xml:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" uri:"end_time" xml:"end_time" form:"end_time"`
	Username  string `json:"username" uri:"username" xml:"username" form:"username"`
	Search    string `json:"search" uri:"search" xml:"search" form:"search"`
	Page      int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset    int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}
