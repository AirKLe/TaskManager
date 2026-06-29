package storage

import (
	"TaskManager/internal/models"
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestNewStorage(t *testing.T) {
	storage := NewInMemoryTaskStorage()
	if storage.data == nil {
		t.Errorf("Expected storage.data is init")
	}
}

func TestGetByID(t *testing.T) {
	storage := NewInMemoryTaskStorage()
	id, err := storage.Create(&models.Task{
		Title:       "NewTask",
		Description: "Napking",
	})
	if err != nil {
		t.Fatal(err)
	}

	newTask, err := storage.GetById(999)
	if err != nil && !errors.Is(err, ErrNotFound) {
		t.Fatal(err)
	}

	newTask, err = storage.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if newTask.Title != "NewTask" {
		t.Errorf("Expected NewTask, got %s", newTask.Title)
	}
	if newTask.Description != "Napking" {
		t.Errorf("Expected TestDesc, got %s", newTask.Description)
	}
}

func TestGetAll(t *testing.T) {
	storage := NewInMemoryTaskStorage()
	tasks := []*models.Task{
		{
			Title:       "NewTask",
			Description: "Napking",
		},
		{
			Title:       "NewTask2",
			Description: "Papking",
		},
		{
			Title:       "NewTask3",
			Description: "Lapking",
		},
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

	for _, task := range newTasks {
		if task.Title != tasks[task.Id-1].Title || task.Description != tasks[task.Id-1].Description {
			t.Errorf("Expected %v, got %v", tasks, newTasks)
		}
	}
}

func TestCreate(t *testing.T) {
	storage := NewInMemoryTaskStorage()
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
		t.Errorf("Expected task should be %v, got %v", task, got)
	}
}

func TestUpdate(t *testing.T) {
	storage := NewInMemoryTaskStorage()
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
	storage := NewInMemoryTaskStorage()
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

func TestCreateIsolation(t *testing.T) {
	store := NewInMemoryTaskStorage()
	aTask := models.Task{
		Title:       "TestTask",
		Description: "TestDesc",
	}

	id, err := store.Create(&aTask)
	if err != nil {
		t.Fatal(err)
	}

	aTask.Title = "x"
	aTask.Description = "x"

	bTask, err := store.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if bTask.Title != "TestTask" {
		t.Errorf("Expected TestTask, got %s", bTask.Title)
	}
	if bTask.Description != "TestDesc" {
		t.Errorf("Expected TestDesc, got %s", bTask.Description)
	}
}

func TestGetByIdIsolation(t *testing.T) {
	store := NewInMemoryTaskStorage()
	aTask := models.Task{
		Title:       "TestTask",
		Description: "TestDesc",
	}

	id, err := store.Create(&aTask)
	if err != nil {
		t.Fatal(err)
	}
	bTask, err := store.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	bTask.Title = "x"
	bTask.Description = "x"

	cTask, err := store.GetById(id)
	if err != nil {
		t.Fatal(err)
	}

	if cTask.Title != "TestTask" {
		t.Errorf("Expected TestTask, got %s", cTask.Title)
	}
	if cTask.Description != "TestDesc" {
		t.Errorf("Expected TestDesc, got %s", cTask.Description)
	}
}

func TestGetAllIsolation(t *testing.T) {
	store := NewInMemoryTaskStorage()
	aTask := models.Task{
		Title:       "TestTask",
		Description: "TestDesc",
	}

	id, err := store.Create(&aTask)
	if err != nil {
		t.Fatal(err)
	}
	bTasks, err := store.GetAll()

	bTasks[0].Title = "x"
	bTasks[0].Description = "x"

	cTask, err := store.GetById(id)

	if cTask.Title != "TestTask" {
		t.Errorf("Expected TestTask, got %s", cTask.Title)
	}
	if cTask.Description != "TestDesc" {
		t.Errorf("Expected TestDesc, got %s", cTask.Description)
	}
}

func TestUpdateIsolation(t *testing.T) {
	store := NewInMemoryTaskStorage()
	aTask := models.Task{
		Title:       "TestTask",
		Description: "TestDesc",
	}

	id, err := store.Create(&aTask)
	if err != nil {
		t.Fatal(err)
	}

	bTask := models.Task{
		Id:          id,
		Title:       "NewTestTask",
		Description: "NewTestDesc",
	}

	store.Update(&bTask)

	bTask.Title = "x"
	bTask.Description = "x"

	cTask, _ := store.GetById(bTask.Id)
	if cTask.Title != "NewTestTask" {
		t.Errorf("Expected NewTestTask, got %s", cTask.Title)
	}
	if cTask.Description != "NewTestDesc" {
		t.Errorf("Expected NewTestDesc, got %s", cTask.Description)
	}
}

func TestConcurrentAccess(t *testing.T) {
	storage := NewInMemoryTaskStorage()
	var wg sync.WaitGroup

	const count = 100

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			task := &models.Task{Title: fmt.Sprintf("Task%d", i+1)}
			storage.Create(task)
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
