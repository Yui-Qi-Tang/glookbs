package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"glookbs.github.com/docs"
	"glookbs.github.com/entity"
	"glookbs.github.com/storage"
)

func init() {
	docs.SwaggerInfo.Title = "Task rest api application doc."
	docs.SwaggerInfo.Description = "it's a task rest api application"
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

// New returns http handler which is implemented by go-gin
func New(mode string, storage *storage.Storage) http.Handler {
	task := &Task{
		db: storage,
	}
	r := gin.Default()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	tasks := r.Group("/tasks")
	{
		tasks.GET("", task.Get)
		tasks.POST("", task.Post)
		tasks.PUT("/:id", task.Put)
		tasks.DELETE("/:id", task.Delete)
	}

	return r
}

type Task struct {
	db *storage.Storage
}

// Get returns tasks
// @Summary returns tasks
// @tags tasks
// @Param page query uint false "1"
// @Param page_size query uint false "10"
// @Produce json
// @Success 200 {array} RespTaskPagination
// @Failure 400 {object} RespErr
// @Failure 500 {object} RespErr
// @Router /tasks [get]
func (t *Task) Get(c *gin.Context) {
	var query RequestGetTaskQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, RespErr{Err: err.Error()})
		return
	}
	data := t.db.Range(query.Page, query.PageSize)
	result := RespTaskPagination{
		Total:    t.db.Count(),
		Page:     query.Page,
		PageSize: query.PageSize,
		Tasks:    make([]RespTask, 0, len(data)),
	}
	for i := range data {
		rt := RespTask{
			ID:     data[i].(*entity.Task).ID,
			Name:   data[i].(*entity.Task).Name,
			Status: int(data[i].(*entity.Task).Status),
		}
		result.Tasks = append(result.Tasks, rt)
	}
	c.JSON(http.StatusOK, result)
}

// Post creates a task
// @Summary create task
// @tags tasks
// @Accept  json
// @Param request body RequsetCreateTask true "request data"
// @Produce json
// @Success 200 {object} RespCreateTaskOK
// @Failure 400 {object} RespErr
// @Failure 500 {object} RespErr
// @Router /tasks [post]
func (t *Task) Post(c *gin.Context) {
	var req RequsetCreateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespErr{Err: err.Error()})
		return
	}
	task := entity.Task{
		Name:   req.Name,
		Status: entity.TaskStatus(req.Status),
	}
	id, err := t.db.Insert(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespErr{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, RespCreateTaskOK{ID: id})
}

// Put create or update task by id
// @Summary create or update task by id
// @tags tasks
// @Param id path string true "id"
// @Param request body RequsetCreateTask true "request data"
// @Success 200 {object} RespTask
// @Failure 400 {object} RespErr
// @Failure 500 {object} RespErr
// @Router /tasks/{id} [put]
func (t *Task) Put(c *gin.Context) {
	var req RequestPutTask
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespErr{Err: err.Error()})
		return
	}

	var reqCreate RequsetCreateTask
	if err := c.ShouldBindJSON(&reqCreate); err != nil {
		c.JSON(http.StatusBadRequest, RespErr{Err: err.Error()})
		return
	}
	task := entity.Task{
		ID:     req.ID,
		Name:   reqCreate.Name,
		Status: entity.TaskStatus(reqCreate.Status),
	}
	if err := t.db.Update(req.ID, &task); err == nil {
		c.JSON(http.StatusOK, RespTask{
			ID:     task.ID,
			Name:   task.Name,
			Status: int(task.Status),
		})
		return
	}

	_, err := t.db.Insert(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, RespErr{Err: err.Error()})
		return
	}
	c.JSON(http.StatusOK, RespTask{
		ID:     task.ID,
		Name:   task.Name,
		Status: int(task.Status),
	})
}

// Delete deletes task by id
// @Summary deletes task by id
// @tags tasks
// @Param id path string true "id"
// @Success 202
// @Failure 400 {object} RespErr
// @Failure 500 {object} RespErr
// @Router /tasks/{id} [delete]
func (t *Task) Delete(c *gin.Context) {
	var req RequestDeleteTask
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, RespErr{Err: err.Error()})
		return
	}
	if err := t.db.Delete(req.ID); err != nil {
		c.JSON(http.StatusInternalServerError, RespErr{Err: err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusAccepted)
}
