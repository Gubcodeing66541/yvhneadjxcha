package Request

type Page struct {
	Count       int `json:"count"`
	CurrentPage int `json:"current_page" binding:"required"`
	CurrentSize int `json:"current_size" binding:"required"`
}

type Count struct {
	Count int `json:"count" uri:"count" form:"count" binding:"required"`
}

type Id struct {
	Id int `json:"id" uri:"id" form:"id" binding:"required"`
}

type FileName struct {
	FileName string `json:"file_name" uri:"file_name" form:"file_name" binding:"required"`
}

type SocketMsg struct {
	ServiceId int    `json:"service_id" uri:"service_id" form:"service_id" binding:"required"`
	Type      string `json:"type" uri:"type" form:"type" binding:"required"`
	Content   string `json:"content" uri:"content" form:"content" binding:"required"`
}
