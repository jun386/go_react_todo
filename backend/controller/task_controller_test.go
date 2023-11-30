package controller_test

import (
	"backend/controller"
	"backend/model"
	"backend/usecase/usecase_mocks"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetAllTasks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskUsecase := usecase_mocks.NewMockITaskUsecase(mockCtrl)
	taskController := controller.NewTaskController(mockTaskUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockTasks := []model.TaskResponse{}
	mockTaskUsecase.EXPECT().GetAllTasks().Return(mockTasks, nil)

	if assert.NoError(t, taskController.GetAllTasks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var tasks []model.TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &tasks); assert.NoError(t, err) {
			assert.Equal(t, mockTasks, tasks)
		}
	}
}

func TestGetTaskById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskUsecase := usecase_mocks.NewMockITaskUsecase(mockCtrl)
	taskController := controller.NewTaskController(mockTaskUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("taskId")
	c.SetParamValues("1")

	now := time.Now().UTC()
	expectedTaskResponse := model.TaskResponse{ID: 1, Title: "Sample Task", CreatedAt: now, UpdatedAt: now}
	mockTaskUsecase.EXPECT().GetTaskById(uint(1)).Return(expectedTaskResponse, nil)

	if assert.NoError(t, taskController.GetTaskById(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var task model.TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &task); assert.NoError(t, err) {
			assert.Equal(t, expectedTaskResponse, task)
		}
	}
}

func TestCreateTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskUsecase := usecase_mocks.NewMockITaskUsecase(mockCtrl)
	taskController := controller.NewTaskController(mockTaskUsecase)

	e := echo.New()
	taskJSON := `{"title": "New Task", "description": "This is a new task"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(taskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	now := time.Now().UTC()
	createdTask := model.TaskResponse{ID: 1, Title: "New Task", CreatedAt: now, UpdatedAt: now}
	mockTaskUsecase.EXPECT().CreateTask(gomock.Any()).Return(createdTask, nil)

	if assert.NoError(t, taskController.CreateTask(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var task model.TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &task); assert.NoError(t, err) {
			assert.Equal(t, createdTask, task)
		}
	}
}

func TestUpdateTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskUsecase := usecase_mocks.NewMockITaskUsecase(mockCtrl)
	taskController := controller.NewTaskController(mockTaskUsecase)

	e := echo.New()
	taskJSON := `{"title": "Updated Task", "description": "This task has been updated"}`
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(taskJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("taskId")
	c.SetParamValues("1")

	now := time.Now().UTC()
	updatedTaskWithID := model.TaskResponse{ID: 1, Title: "Updated Task", CreatedAt: now, UpdatedAt: now}
	mockTaskUsecase.EXPECT().UpdateTask(gomock.Any(), uint(1)).Return(updatedTaskWithID, nil)

	if assert.NoError(t, taskController.UpdateTask(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var task model.TaskResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &task); assert.NoError(t, err) {
			assert.Equal(t, updatedTaskWithID, task)
		}
	}
}

func TestDeleteTask(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockTaskUsecase := usecase_mocks.NewMockITaskUsecase(mockCtrl)
	taskController := controller.NewTaskController(mockTaskUsecase)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("taskId")
	c.SetParamValues("1")

	mockTaskUsecase.EXPECT().DeleteTask(uint(1)).Return(nil)

	if assert.NoError(t, taskController.DeleteTask(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
