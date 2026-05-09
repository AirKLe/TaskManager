package storage

import (
	"TaskManager/internal/models"
	"errors"
	"fmt"
	"sync"
)

type TaskStorage interface {
	GetById(id int) (*models.Task, error)
	GetAll() ([]*models.Task, error)
	Create(t *models.Task) (int, error)
	Update(t *models.Task) error
	Delete(id int) error
}

var ErrNotFound = errors.New("not found")

type inMemoryTaskStorage struct {
	mu     sync.RWMutex
	data   map[int]*models.Task
	nextId int
}

func NewInMemoryTaskStorage() *inMemoryTaskStorage {
	return &inMemoryTaskStorage{
		mu:     sync.RWMutex{},
		data:   make(map[int]*models.Task),
		nextId: 1,
	}
}

func (s *inMemoryTaskStorage) GetById(id int) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.data[id]
	if !ok {
		return nil, fmt.Errorf("task %v: %w", id, ErrNotFound)
	}

	return task.Clone(), nil
}

func (s *inMemoryTaskStorage) GetAll() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.data))
	for _, t := range s.data {
		tasks = append(tasks, t.Clone())
	}
	return tasks, nil
}

func (s *inMemoryTaskStorage) Create(t *models.Task) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copy := t.Clone()
	copy.Id = s.nextId
	s.data[copy.Id] = copy
	s.nextId++
	return copy.Id, nil
}

func (s *inMemoryTaskStorage) Update(t *models.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[t.Id]; !ok {
		return fmt.Errorf("task %v: %w", t.Id, ErrNotFound)
	}
	s.data[t.Id] = t.Clone()
	return nil
}

func (s *inMemoryTaskStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; ok {
		delete(s.data, id)
		return nil
	}

	return fmt.Errorf("task %v: %w", id, ErrNotFound)
}
