package service

import (
	"TaskManager/internal/mocks"
	"TaskManager/internal/models"
	"TaskManager/internal/storage"
	"context"
	"errors"
	"testing"
)

func TestGetTask_Success(t *testing.T) {
	expected := &models.Task{
		Id:          36,
		Title:       "NewTask",
		Description: "Napking",
	}

	mockStore := &mocks.MockTaskStorage{
		GetByIdFunc: func(ctx context.Context, id int) (*models.Task, error) {
			if id != 36 {
				t.Errorf("Expected id 36, got %d", id)
			}
			return expected, nil
		},
	}
	service := NewTaskService(mockStore)

	got, err := service.GetTask(context.Background(), 36)
	if err != nil {
		t.Fatal(err)
	}

	if *expected != *got {
		t.Errorf("Expected %v, got %v", *expected, *got)
	}
}

func TestGetTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	_, err := service.GetTask(context.Background(), -2)

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestCreateTask_Success(t *testing.T) {
	mockStore := &mocks.MockTaskStorage{
		CreateFunc: func(ctx context.Context, t *models.Task) (int, error) {
			return t.Id, nil
		},
	}

	service := NewTaskService(mockStore)

	id, err := service.CreateTask(
		context.Background(),
		&models.Task{
			Id:          36,
			Title:       "NewTask",
			Description: "Napking",
		})

	if err != nil {
		t.Fatal(err)
	}

	if id != 36 {
		t.Errorf("Expected 36, got %v", id)
	}
}

func TestCreateTask_NilTask(t *testing.T) {
	service := NewTaskService(nil)

	_, err := service.CreateTask(context.Background(), nil)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateTask_Success(test *testing.T) {
	expected := &models.Task{
		Id:    36,
		Title: "NewTask",
	}

	mockStore := &mocks.MockTaskStorage{
		UpdateFunc: func(ctx context.Context, t *models.Task) error {
			if t.Id != 36 {
				test.Errorf("Expected 36, got %v", t.Id)
			}
			return nil
		},
	}

	service := NewTaskService(mockStore)

	err := service.UpdateTask(context.Background(), expected)

	if err != nil {
		test.Fatal(err)
	}
}

func TestUpdateTask_NilTask(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(context.Background(), nil)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(
		context.Background(),
		&models.Task{
			Id: -2,
		})

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestUpdateTask_EmptyTitle(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(context.Background(), &models.Task{})

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockStore := &mocks.MockTaskStorage{
		UpdateFunc: func(ctx context.Context, t *models.Task) error {
			return storage.ErrNotFound
		},
	}

	service := NewTaskService(mockStore)

	err := service.UpdateTask(
		context.Background(),
		&models.Task{
			Id:          36,
			Title:       "NewTask",
			Description: "Napking",
		})

	var notFoundErr *NotFoundError

	if !errors.As(err, &notFoundErr) {
		t.Fatalf("Expected notFoundError, got %v", err)
	}
}

func TestDeleteTask_Success(test *testing.T) {
	expectedId := 36

	mockStore := &mocks.MockTaskStorage{
		DeleteFunc: func(ctx context.Context, id int) error {
			if id != expectedId {
				test.Errorf("Expected 36, got %v", id)
			}
			return nil
		},
	}

	service := NewTaskService(mockStore)

	err := service.DeleteTask(context.Background(), expectedId)

	if err != nil {
		test.Fatal(err)
	}
}

func TestDeleteTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	err := service.DeleteTask(context.Background(), -2)

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}
