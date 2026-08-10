package storage

import (
	"TaskManager/internal/config"
	"TaskManager/internal/models"
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestStorage(t *testing.T) *PostgresTaskStorage {
	t.Helper()

	cfg, err := config.Load("../../config_test.json")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DataBaseURL())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(pool.Close)

	return NewPostgresTaskStorage(pool)
}

func clearTasks(t *testing.T, storage *PostgresTaskStorage) {
	t.Helper()

	_, err := storage.db.Exec(
		context.Background(),
		"TRUNCATE TABLE tasks RESTART IDENTITY",
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewStorage(t *testing.T) {
	storage := newTestStorage(t)

	if storage.db == nil {
		t.Error("Expected storage.db to be initialized")
	}
}

func TestGetByID(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	id, err := storage.Create(&models.Task{
		Title:       "NewTask",
		Description: "Napking",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = storage.GetById(999999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Expected ErrNotFound, got %v", err)
	}

	newTask, err := storage.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if newTask.Title != "NewTask" {
		t.Errorf("Expected NewTask, got %s", newTask.Title)
	}

	if newTask.Description != "Napking" {
		t.Errorf("Expected Napking, got %s", newTask.Description)
	}
}

func TestGetAll(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	tasks := []*models.Task{
		{Title: "NewTask", Description: "Napking"},
		{Title: "NewTask2", Description: "Papking"},
		{Title: "NewTask3", Description: "Lapking"},
	}

	for _, task := range tasks {
		_, err := storage.Create(task)
		if err != nil {
			t.Fatal(err)
		}
	}

	newTasks, err := storage.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(newTasks) < len(tasks) {
		t.Fatalf("Expected at least %d tasks, got %d", len(tasks), len(newTasks))
	}
}

func TestCreate(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	task := &models.Task{
		Title:       "NewTask",
		Description: "Napking",
	}

	id, err := storage.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != task.Title || got.Description != task.Description {
		t.Errorf("Expected task %v, got %v", task, got)
	}
}

func TestUpdate(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	task := &models.Task{
		Title:       "NewTask",
		Description: "Napking",
	}
	id, err := storage.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	newTask := &models.Task{
		Id:          id,
		Title:       "New title",
		Description: "New desc",
	}
	err = storage.Update(newTask)
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != newTask.Title || got.Description != newTask.Description {
		t.Errorf("Expected task should be %v, got %v", newTask, got)
	}
}

func TestDelete(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	task := &models.Task{
		Title:       "NewTask",
		Description: "Napking",
	}
	id, err := storage.Create(task)
	if err != nil {
		t.Fatal(err)
	}

	err = storage.Delete(id)
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.GetById(id)
	if got != nil && !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected task has not deleted:%v", task)
	}
}

func TestConcurrentAccess(t *testing.T) {
	storage := newTestStorage(t)
	clearTasks(t, storage)

	var wg sync.WaitGroup

	const count = 100

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			task := &models.Task{
				Title: fmt.Sprintf("Task%d", i+1),
			}

			_, err := storage.Create(task)
			if err != nil {
				t.Errorf("Create failed: %v", err)
			}
		}(i)
	}

	wg.Wait()

	newTasks, err := storage.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(newTasks) != count {
		t.Fatalf("Expected %d, got %d", count, len(newTasks))
	}
}
