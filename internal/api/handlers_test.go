package api

import (
	"TaskManager/internal/mocks"
	"TaskManager/internal/models"
	"TaskManager/internal/service"
	"encoding/json"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetTask_Success(t *testing.T) {
	expected := &models.Task{
		Id:          36,
		Title:       "NewTask",
		Description: "Napking",
	}

	mockStore := &mocks.MockTaskStorage{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return expected, nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/task?id=36",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestHandleGetTask_InvalidId(t *testing.T) {
	expected := &models.Task{
		Id:          -1,
		Title:       "NewTask",
		Description: "Napking",
	}

	mockStore := &mocks.MockTaskStorage{
		GetByIdFunc: func(id int) (*models.Task, error) {
			return expected, nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/task?id=-1",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected 400, got %d", resp.StatusCode)
	}
}

func TestHandleCreateTask_Success(t *testing.T) {
	mockStore := &mocks.MockTaskStorage{
		CreateFunc: func(t *models.Task) (int, error) {
			return t.Id, nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	body := `{
		"title": "NewTask",
		"description": "Napking"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/task",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected 201 , got %d", w.Code)
	}
}

func TestHandleCreateTask_InvalidJSON(t *testing.T) {
	mockStore := &mocks.MockTaskStorage{
		CreateFunc: func(t *models.Task) (int, error) {
			return t.Id, nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	body := `{
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/task",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 , got %d", w.Code)
	}
}

func TestHandleUpdateTask_Success(test *testing.T) {
	mockStore := &mocks.MockTaskStorage{
		UpdateFunc: func(t *models.Task) error {
			if t.Id != 36 {
				test.Errorf("Expected 36, got %v", t.Id)
			}
			return nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	body := `{
		"title": "NewTask",
		"description": "Napking"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/task?id=36",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		test.Fatalf("Expected 200 , got %d", w.Code)
	}

	got := models.Task{}

	err := json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		test.Fatal(err)
	}

	if got.Id != 36 {
		test.Fatalf("Expected id 36, got %d", got.Id)
	}
}

func TestHandleUpdateTask_InvalidJSON(t *testing.T) {
	srvc := service.NewTaskService(nil)

	handler := NewTaskHandler(srvc)

	body := `{
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/task?id=36",
		strings.NewReader(body),
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 , got %d", w.Code)
	}
}

func TestHandleDeleteTask_Success(test *testing.T) {
	expectedId := 36

	mockStore := &mocks.MockTaskStorage{
		DeleteFunc: func(id int) error {
			if id != expectedId {
				test.Errorf("Expected 36, got %v", id)
			}
			return nil
		},
	}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/task?id=36",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		test.Fatalf("Expected 204 , got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mockStore := &mocks.MockTaskStorage{}

	srvc := service.NewTaskService(mockStore)

	handler := NewTaskHandler(srvc)

	req := httptest.NewRequest(
		http.MethodPatch,
		"/task",
		nil,
	)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected 405 , got %d", w.Code)
	}
}
