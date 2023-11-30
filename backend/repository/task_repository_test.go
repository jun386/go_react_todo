package repository

import (
	"backend/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGetAllTasks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("スタブデータベース接続を開く際にエラー '%s' が発生しました", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gormデータベースを開く際にエラー '%s' が発生しました", err)
	}

	tasks := []model.Task{{}, {}}
	mock.ExpectQuery("^SELECT (.+) FROM \"tasks\" ORDER BY created_at$").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "テストタスク1").AddRow(2, "テストタスク2"))

	repo := NewTaskRepository(gormDB)
	err = repo.GetAllTasks(&tasks)

	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	assert.Equal(t, "テストタスク1", tasks[0].Title)
	assert.Equal(t, "テストタスク2", tasks[1].Title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未達の期待値があります: %s", err)
	}
}

func TestGetTaskById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("スタブデータベース接続を開く際にエラー '%s' が発生しました", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gormデータベースを開く際にエラー '%s' が発生しました", err)
	}

	taskId := uint(1)
	mock.ExpectQuery("^SELECT (.+) FROM \"tasks\" WHERE (.+)").
		WithArgs(taskId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(taskId, "テストタスク"))

	repo := NewTaskRepository(gormDB)
	task := model.Task{}
	err = repo.GetTaskById(&task, taskId)

	assert.NoError(t, err)
	assert.Equal(t, taskId, task.ID)
	assert.Equal(t, "テストタスク", task.Title)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未達の期待値があります: %s", err)
	}
}

func TestCreateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("スタブデータベース接続を開く際にエラー '%s' が発生しました", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gormデータベースを開く際にエラー '%s' が発生しました", err)
	}

	task := &model.Task{Title: "新しいタスク"}
	mock.ExpectBegin()
	mock.ExpectQuery("^INSERT INTO \"tasks\"").
		WithArgs(task.Title, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	repo := NewTaskRepository(gormDB)
	err = repo.CreateTask(task)

	assert.NoError(t, err)
	assert.NotZero(t, task.ID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未達の期待値があります: %s", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("スタブデータベース接続を開く際にエラー '%s' が発生しました", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gormデータベースを開く際にエラー '%s' が発生しました", err)
	}

	taskId := uint(1)
	updatedTask := &model.Task{ID: taskId, Title: "更新されたタスク"}
	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectQuery("^UPDATE \"tasks\"").
		WithArgs(updatedTask.Title, sqlmock.AnyArg(), taskId, taskId).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "created_at", "updated_at"}).AddRow(taskId, updatedTask.Title, now, now))
	mock.ExpectCommit()

	repo := NewTaskRepository(gormDB)
	err = repo.UpdateTask(updatedTask, taskId)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未達の期待値があります: %s", err)
	}
}

func TestDeleteTask(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("スタブデータベース接続を開く際にエラー '%s' が発生しました", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gormデータベースを開く際にエラー '%s' が発生しました", err)
	}

	taskId := uint(1)
	mock.ExpectBegin()
	mock.ExpectExec("^DELETE FROM \"tasks\"").
		WithArgs(taskId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewTaskRepository(gormDB)
	err = repo.DeleteTask(taskId)

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("未達の期待値があります: %s", err)
	}
}
