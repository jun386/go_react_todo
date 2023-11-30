package usecase

import (
	"backend/model"
	"backend/repository/repository_mocks"
	"backend/validator/validator_mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetAllTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := repository_mocks.NewMockITaskRepository(ctrl)
	mockTaskValidator := validator_mocks.NewMockITaskValidator(ctrl)
	u := NewTaskUsecase(mockTaskRepo, mockTaskValidator)

	mockTasks := []model.Task{
		{ID: 1, Title: "Task 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Title: "Task 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	mockTaskRepo.EXPECT().GetAllTasks(gomock.Any()).DoAndReturn(func(tasks *[]model.Task) error {
		*tasks = mockTasks
		return nil
	})

	res, err := u.GetAllTasks()
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Len(t, res, 2)
}

func TestGetTaskById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := repository_mocks.NewMockITaskRepository(ctrl)
	mockTaskValidator := validator_mocks.NewMockITaskValidator(ctrl)
	u := NewTaskUsecase(mockTaskRepo, mockTaskValidator)

	taskID := uint(1)
	mockTask := model.Task{
		ID:        taskID,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Success", func(t *testing.T) {
		mockTaskRepo.EXPECT().GetTaskById(gomock.Any(), taskID).DoAndReturn(func(task *model.Task, id uint) error {
			*task = mockTask
			return nil
		})

		res, err := u.GetTaskById(taskID)
		assert.NoError(t, err)
		assert.Equal(t, mockTask.ID, res.ID)
		assert.Equal(t, mockTask.Title, res.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		notFoundID := uint(2)
		mockTaskRepo.EXPECT().GetTaskById(gomock.Any(), notFoundID).Return(errors.New("not found"))

		_, err := u.GetTaskById(notFoundID)
		assert.Error(t, err)
	})
}

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := repository_mocks.NewMockITaskRepository(ctrl)
	mockTaskValidator := validator_mocks.NewMockITaskValidator(ctrl)
	u := NewTaskUsecase(mockTaskRepo, mockTaskValidator)

	newTask := model.Task{
		Title: "New Task",
	}

	t.Run("Success", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(newTask).Return(nil)
		mockTaskRepo.EXPECT().CreateTask(gomock.Any()).DoAndReturn(func(task *model.Task) error {
			task.ID = 1 // Assuming the task is assigned an ID of 1 upon creation
			task.CreatedAt = time.Now()
			task.UpdatedAt = time.Now()
			return nil
		})

		res, err := u.CreateTask(newTask)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), res.ID)
		assert.Equal(t, newTask.Title, res.Title)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(newTask).Return(errors.New("validation error"))

		_, err := u.CreateTask(newTask)
		assert.Error(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(newTask).Return(nil)
		mockTaskRepo.EXPECT().CreateTask(gomock.Any()).Return(errors.New("repository error"))

		_, err := u.CreateTask(newTask)
		assert.Error(t, err)
	})
}

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := repository_mocks.NewMockITaskRepository(ctrl)
	mockTaskValidator := validator_mocks.NewMockITaskValidator(ctrl)
	u := NewTaskUsecase(mockTaskRepo, mockTaskValidator)

	taskID := uint(1)
	updatedTask := model.Task{
		ID:    taskID,
		Title: "Updated Task",
	}

	t.Run("Success", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(updatedTask).Return(nil)
		mockTaskRepo.EXPECT().UpdateTask(gomock.Any(), taskID).DoAndReturn(func(task *model.Task, id uint) error {
			task.UpdatedAt = time.Now()
			return nil
		})

		res, err := u.UpdateTask(updatedTask, taskID)
		assert.NoError(t, err)
		assert.Equal(t, updatedTask.ID, res.ID)
		assert.Equal(t, updatedTask.Title, res.Title)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(updatedTask).Return(errors.New("validation error"))

		_, err := u.UpdateTask(updatedTask, taskID)
		assert.Error(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockTaskValidator.EXPECT().TaskValidate(updatedTask).Return(nil)
		mockTaskRepo.EXPECT().UpdateTask(gomock.Any(), taskID).Return(errors.New("repository error"))

		_, err := u.UpdateTask(updatedTask, taskID)
		assert.Error(t, err)
	})
}

func TestDeleteTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := repository_mocks.NewMockITaskRepository(ctrl)
	u := NewTaskUsecase(mockTaskRepo, nil) // バリデーターはDeleteTaskには不要です。

	taskID := uint(1)

	t.Run("Success", func(t *testing.T) {
		mockTaskRepo.EXPECT().DeleteTask(taskID).Return(nil)

		err := u.DeleteTask(taskID)
		assert.NoError(t, err)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockTaskRepo.EXPECT().DeleteTask(taskID).Return(errors.New("repository error"))

		err := u.DeleteTask(taskID)
		assert.Error(t, err)
	})
}
