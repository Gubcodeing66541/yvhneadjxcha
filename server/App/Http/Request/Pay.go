package Request

type PayRecorder struct {
	Reason           string `json:"reason" uri:"reason" xml:"reason" form:"reason"`
	Search           string `json:"search" uri:"search" xml:"search" form:"search"`
	ServiceManagerId int    `json:"service_manager_id"  uri:"service_manager_id" xml:"service_manager_id" form:"service_manager_id"`
	StartTime        string `json:"start_time" uri:"start_time" xml:"start_time" form:"start_time"`
	EndTime          string `json:"end_time" uri:"end_time" xml:"end_time" form:"end_time"`
	Page             int    `json:"page" uri:"page" xml:"page" form:"page"`
	Offset           int    `json:"offset" uri:"offset" xml:"offset" form:"offset"`
}
