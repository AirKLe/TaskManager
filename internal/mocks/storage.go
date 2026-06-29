package mocks

import "TaskManager/internal/models"

type MockTaskStorage struct {
	GetByIdFunc func(id int) (*models.Task, error)
	GetAllFunc  func() ([]*models.Task, error)
	CreateFunc  func(t *models.Task) (int, error)
	UpdateFunc  func(t *models.Task) error
	DeleteFunc  func(id int) error
}

func (m *MockTaskStorage) GetById(id int) (*models.Task, error) {
	return m.GetByIdFunc(id)
}

func (m *MockTaskStorage) GetAll() ([]*models.Task, error) {
	return m.GetAllFunc()
}

func (m *MockTaskStorage) Create(t *models.Task) (int, error) {
	return m.CreateFunc(t)
}

func (m *MockTaskStorage) Update(t *models.Task) error {
	return m.UpdateFunc(t)
}

func (m *MockTaskStorage) Delete(id int) error {
	return m.DeleteFunc(id)
}
