package httphandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"glookbs.github.com/storage"
	"glookbs.github.com/storage/drivers/skiplists"
	"gotest.tools/assert"
)

type Data struct {
	ID   int
	Name string
}

func TestCreateTask(t *testing.T) {
	router := New(gin.TestMode, storage.New(skiplists.New()))
	testcase := []struct {
		requests []RequsetCreateTask
		want     []RespCreateTaskOK
	}{
		{
			requests: []RequsetCreateTask{
				{
					Name:   "t1",
					Status: 0,
				},
				{
					Name:   "t2",
					Status: 1,
				},
				{
					Name:   "t3",
					Status: 1,
				},
			},
			want: []RespCreateTaskOK{{1}, {2}, {3}},
		},
	}

	for _, tt := range testcase {
		for i := range tt.requests {
			data, err := json.Marshal(tt.requests[i])
			if err != nil {
				t.Fatalf("json marshal error: %v", err)
			}
			req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(data))
			if err != nil {
				t.Fatalf("create http request error %v", err)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			// check response
			var resp RespCreateTaskOK
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("json unmarshal error: %v", err)
			}
			assert.Equal(t, tt.want[i], resp)
		}
	}
}

func TestGetTasks(t *testing.T) {
	router := New(gin.TestMode, storage.New(skiplists.New()))
	// add tasks
	requests := []RequsetCreateTask{
		{
			Name:   "t1",
			Status: 0,
		},
		{
			Name:   "t2",
			Status: 1,
		},
		{
			Name:   "t3",
			Status: 1,
		},
	}

	expected := RespTaskPagination{
		Total:    len(requests),
		Page:     1,
		PageSize: 10,
		Tasks:    make([]RespTask, 0, len(requests)),
	}

	for i, req := range requests {
		expected.Tasks = append(expected.Tasks, RespTask{
			ID:     i + 1,
			Name:   req.Name,
			Status: req.Status,
		})
	}

	for i := range requests {
		data, err := json.Marshal(requests[i])
		if err != nil {
			t.Fatalf("json marshal error: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("create http request error %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// get tasks
	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp RespTaskPagination
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal resp: %v", err)
	}
	assert.DeepEqual(t, resp, expected)
}

func TestDeleteTask(t *testing.T) {
	router := New(gin.TestMode, storage.New(skiplists.New()))
	// add tasks
	requests := []RequsetCreateTask{
		{
			Name:   "t1",
			Status: 0,
		},
		{
			Name:   "t2",
			Status: 1,
		},
		{
			Name:   "t3",
			Status: 1,
		},
	}

	expected := RespTaskPagination{
		Total:    2,
		Page:     1,
		PageSize: 10,
		Tasks: []RespTask{
			{
				ID:     1,
				Name:   "t1",
				Status: 0,
			},
			{
				ID:     3,
				Name:   "t3",
				Status: 1,
			},
		},
	}

	for i := range requests {
		data, err := json.Marshal(requests[i])
		if err != nil {
			t.Fatalf("json marshal error: %v", err)
		}
		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("create http request error %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// delete tasks
	req, err := http.NewRequest(http.MethodDelete, "/tasks/2", nil)
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusAccepted, w.Code)

	// get tasks
	req, err = http.NewRequest(http.MethodGet, "/tasks", nil)
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp RespTaskPagination
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal resp: %v", err)
	}
	assert.DeepEqual(t, resp, expected)
}

func TestCreateOrUpdateTask(t *testing.T) {
	router := New(gin.TestMode, storage.New(skiplists.New()))
	// add tasks
	requests := []RequsetCreateTask{
		{
			Name:   "t1",
			Status: 0,
		},
		{
			Name:   "t2",
			Status: 1,
		},
		{
			Name:   "t3",
			Status: 1,
		},
	}

	expected := RespTaskPagination{
		Total:    3,
		Page:     1,
		PageSize: 10,
		Tasks: []RespTask{
			{
				ID:     1,
				Name:   "t1",
				Status: 0,
			},
			{
				ID:     2,
				Name:   "t2",
				Status: 1,
			},
			{
				ID:     3,
				Name:   "t3",
				Status: 1,
			},
		},
	}

	// create data by put
	for i := range requests {
		data, err := json.Marshal(requests[i])
		if err != nil {
			t.Fatalf("json marshal error: %v", err)
		}
		uri := fmt.Sprintf("/tasks/%d", i+1)
		req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(data))
		if err != nil {
			t.Fatalf("create http request error %v", err)
		}
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	// get tasks
	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp RespTaskPagination
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal resp: %v", err)
	}
	assert.DeepEqual(t, resp, expected)

	// modify data with 1
	newTask := RequsetCreateTask{Name: "rename-t1", Status: 1}
	data, err := json.Marshal(newTask)
	if err != nil {
		t.Fatalf("json marshal error: %v", err)
	}
	req, err = http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var respTaskUpdated RespTask
	if err := json.Unmarshal(w.Body.Bytes(), &respTaskUpdated); err != nil {
		t.Fatalf("failed to unmarshal resp: %v", err)
	}
	expectedUpdateTask := RespTask{
		ID:     1,
		Name:   "rename-t1",
		Status: 1,
	}
	assert.DeepEqual(t, respTaskUpdated, expectedUpdateTask)
}

func TestMultipleClientsCreateTaskSimultaneously(t *testing.T) {
	router := New(gin.TestMode, storage.New(skiplists.New()))
	clients := 10
	requests := make([]RequsetCreateTask, 0, clients)

	for i := 0; i < clients; i++ {
		requests = append(requests, RequsetCreateTask{
			Name:   fmt.Sprintf("client-%d-task", i),
			Status: i % 2,
		})
	}

	var wg sync.WaitGroup
	wg.Add(clients)
	for i := 0; i < clients; i++ {
		go func(task RequsetCreateTask) {
			defer wg.Done()
			data, _ := json.Marshal(task)
			req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(data))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}(requests[i])
	}
	wg.Wait()

	// get tasks
	req, err := http.NewRequest(http.MethodGet, "/tasks", nil)
	if err != nil {
		t.Fatalf("create http request error %v", err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp RespTaskPagination
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal resp: %v", err)
	}
	assert.Equal(t, len(resp.Tasks), len(requests))

	// check if id is unique
	exist := make(map[int]bool, len(resp.Tasks))
	for _, task := range resp.Tasks {
		if exist[task.ID] {
			t.Fatal("duplicated id", task.ID)
		}
		exist[task.ID] = true
	}
}
