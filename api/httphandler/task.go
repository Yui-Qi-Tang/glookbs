package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"glookbs.github.com/storage"
)

// New returns http handler which is implemented by go-gin
func New(mode string, storage *storage.Storage) http.Handler {
	task := &Task{
		db: storage,
	}
	r := gin.Default()

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

func (t *Task) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
func (t *Task) Post(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
func (t *Task) Put(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
func (t *Task) Delete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
}
