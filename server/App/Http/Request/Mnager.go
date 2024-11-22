package Request

type Account struct {
	ServiceManagerId int `json:"service_manager_id" uri:"service_manager_id" form:"service_manager_id" binding:"required"`
	Account          int `json:"account" uri:"account" form:"account" binding:"required"`
}
