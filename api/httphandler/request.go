package httphandler

type RequsetCreateTask struct {
	Name   string `json:"name" binding:"required" example:"task-1"`
	Status int    `json:"status" binding:"required,min=0,max=1"`
}

type RequestGetTaskQuery struct {
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1"`
}

type RequestDeleteTask struct {
	ID int `uri:"id" binding:"required,min=1"`
}

type RequestPutTask struct {
	ID int `uri:"id" binding:"required,min=1"`
}
