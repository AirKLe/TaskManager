package service

import (
	"TaskManager/internal/models"
	"TaskManager/internal/storage"
	"errors"
	"testing"
)

type mockTaskStorage struct {
	GetByIdFunc func(id int) (*models.Task, error)
	GetAllFunc  func() ([]*models.Task, error)
	CreateFunc  func(t *models.Task) (int, error)
	UpdateFunc  func(t *models.Task) error
	DeleteFunc  func(id int) error
}

func (m *mockTaskStorage) GetById(id int) (*models.Task, error) {
	return m.GetByIdFunc(id)
}

func (m *mockTaskStorage) GetAll() ([]*models.Task, error) {
	return m.GetAllFunc()
}

func (m *mockTaskStorage) Create(t *models.Task) (int, error) {
	return m.CreateFunc(t)
}

func (m *mockTaskStorage) Update(t *models.Task) error {
	return m.UpdateFunc(t)
}

func (m *mockTaskStorage) Delete(id int) error {
	return m.DeleteFunc(id)
}

func TestGetTask_Success(t *testing.T) {
	expected := &models.Task{
		Id:          36,
		Title:       "NewTask",
		Description: "Napking",
	}

	mockStore := &mockTaskStorage{
		GetByIdFunc: func(id int) (*models.Task, error) {
			if id != 36 {
				t.Errorf("Expected id 36, got %d", id)
			}
			return expected, nil
		},
	}
	service := NewTaskService(mockStore)

	got, err := service.GetTask(36)
	if err != nil {
		t.Fatal(err)
	}

	if *expected != *got {
		t.Errorf("Expected %v, got %v", *expected, *got)
	}
}

func TestGetTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	_, err := service.GetTask(-2)

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestCreateTask_Success(t *testing.T) {
	mockStore := &mockTaskStorage{
		CreateFunc: func(t *models.Task) (int, error) {
			return t.Id, nil
		},
	}

	service := NewTaskService(mockStore)

	id, err := service.CreateTask(&models.Task{
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

	_, err := service.CreateTask(nil)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateTask_Success(test *testing.T) {
	expected := &models.Task{
		Id:    36,
		Title: "NewTask",
	}

	mockStore := &mockTaskStorage{
		UpdateFunc: func(t *models.Task) error {
			if t.Id != 36 {
				test.Errorf("Expected 36, got %v", t.Id)
			}
			return nil
		},
	}

	service := NewTaskService(mockStore)

	err := service.UpdateTask(expected)

	if err != nil {
		test.Fatal(err)
	}
}

func TestUpdateTask_NilTask(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(nil)

	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestUpdateTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(&models.Task{
		Id: -2,
	})

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestUpdateTask_EmptyTitle(t *testing.T) {
	service := NewTaskService(nil)

	err := service.UpdateTask(&models.Task{})

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	mockStore := &mockTaskStorage{
		UpdateFunc: func(t *models.Task) error {
			return storage.ErrNotFound
		},
	}

	service := NewTaskService(mockStore)

	err := service.UpdateTask(&models.Task{
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

	mockStore := &mockTaskStorage{
		DeleteFunc: func(id int) error {
			if id != 36 {
				test.Errorf("Expected 36, got %v", id)
			}
			return nil
		},
	}

	service := NewTaskService(mockStore)

	err := service.DeleteTask(expectedId)

	if err != nil {
		test.Fatal(err)
	}
}

func TestDeleteTask_InvalidId(t *testing.T) {
	service := NewTaskService(nil)

	err := service.DeleteTask(-2)

	var validationErr *ValidationError

	if !errors.As(err, &validationErr) {
		t.Fatalf("Expected ValidationError, got %v", err)
	}
}
