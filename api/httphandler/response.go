package httphandler

type RespErr struct {
	Err string `json:"error"`
}

type RespCreateTaskOK struct {
	ID int `json:"id"`
}

type RespTask struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Status int    `json:"status"`
}

type RespTaskPagination struct {
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
	Total    int        `json:"total"`
	Tasks    []RespTask `json:"tasks"`
}
