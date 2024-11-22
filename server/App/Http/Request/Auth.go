package Request

type RegisterAndLogin struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
}

type ActivateLogin struct {
	Username string `form:"username" json:"username" uri:"username" xml:"username" binding:"required"`
}

type UpdatePassword struct {
	Username    string `form:"username" json:"username" uri:"username" xml:"username" `
	Password    string `form:"password" json:"password" uri:"password" xml:"password" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" uri:"new_password" xml:"new_password" binding:"required"`
}
