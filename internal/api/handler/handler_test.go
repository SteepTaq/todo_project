package handler

import (
	"bytes"
	contex "context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SteepTaq/todo_project/internal/api/config"
	"github.com/SteepTaq/todo_project/internal/api/domain"
	"github.com/SteepTaq/todo_project/internal/api/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)



type createTaskService interface {
	CreateTask(ctx contex.Context, title, description string) (*domain.Task, error)
	GetAllTasks(ctx contex.Context) ([]domain.Task, error)
	GetTaskById(ctx contex.Context, id string) (*domain.Task, error)
	UpdateTask(ctx contex.Context, id, title, description, status string) (*domain.Task, error)
	DeleteTask(ctx contex.Context, id int32) error
	Close()
}

type mockService struct{}

func (m *mockService) CreateTask(ctx contex.Context, title, description string) (*domain.Task, error) {
	return &domain.Task{
		ID:          "1",
		Title:       title,
		Description: description,
		Status:      "pending",
	}, nil
}

func (m *mockService) GetAllTasks(ctx contex.Context) ([]domain.Task, error) {
	return nil, nil
}

func (m *mockService) GetTaskById(ctx contex.Context, id string) (*domain.Task, error) {
	return nil, nil
}

func (m *mockService) UpdateTask(ctx contex.Context, id, title, description, status string) (*domain.Task, error) {
	return nil, nil
}

func (m *mockService) DeleteTask(ctx contex.Context, id int32) error {
	return nil
}

func (m *mockService) Close() {}

func newTestTodoHandler(cfg *config.Config, service createTaskService, producer *kafka.Producer) *TodoHandler {
	return &TodoHandler{
		cfg:      cfg,
		service:  service,
		producer: producer,
	}
}

func TestCreateTask(t *testing.T) {
	cfg := &config.Config{}
	service := &mockService{}
	var producer *kafka.Producer = nil 
	h := newTestTodoHandler(cfg, service, producer)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	body := map[string]string{
		"title":       "Test Task",
		"description": "Test Description",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/create", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp domain.Task
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Test Task", resp.Title)
	assert.Equal(t, "Test Description", resp.Description)
}
