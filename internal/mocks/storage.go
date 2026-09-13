package mocks

import (
	"TaskManager/internal/models"
	"context"
)

type MockTaskStorage struct {
	GetByIdFunc func(ctx context.Context, id int) (*models.Task, error)
	GetAllFunc  func(ctx context.Context) ([]*models.Task, error)
	CreateFunc  func(ctx context.Context, t *models.Task) (int, error)
	UpdateFunc  func(ctx context.Context, t *models.Task) error
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *MockTaskStorage) GetById(ctx context.Context, id int) (*models.Task, error) {
	return m.GetByIdFunc(ctx, id)
}

func (m *MockTaskStorage) GetAll(ctx context.Context) ([]*models.Task, error) {
	return m.GetAllFunc(ctx)
}

func (m *MockTaskStorage) Create(ctx context.Context, t *models.Task) (int, error) {
	return m.CreateFunc(ctx, t)
}

func (m *MockTaskStorage) Update(ctx context.Context, t *models.Task) error {
	return m.UpdateFunc(ctx, t)
}

func (m *MockTaskStorage) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}
