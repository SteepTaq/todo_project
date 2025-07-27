package handler

import (
	contex "context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/SteepTaq/todo_project/internal/api/config"
	"github.com/SteepTaq/todo_project/internal/api/domain"
	"github.com/SteepTaq/todo_project/internal/api/kafka"
	"github.com/SteepTaq/todo_project/pkg/context"
	"github.com/SteepTaq/todo_project/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type TodoHandler struct {
	cfg      *config.Config
	service  DBClientInterface
	producer *kafka.Producer
}
type DBClientInterface interface {
	CreateTask(ctx contex.Context, title, description string) (*domain.Task, error)
	GetAllTasks(ctx contex.Context) ([]domain.Task, error)
	GetTaskById(ctx contex.Context, id string) (*domain.Task, error)
	UpdateTask(ctx contex.Context, id, title, description, status string) (*domain.Task, error)
	DeleteTask(ctx contex.Context, id int32) error
	Close()
}

func NewTodoHandler(cfg *config.Config, service DBClientInterface, producer *kafka.Producer) *TodoHandler {
	return &TodoHandler{
		cfg:      cfg,
		service:  service,
		producer: producer,
	}
}

func (h *TodoHandler) RegisterRoutes(router chi.Router) {
	router.Get("/list", h.GetAllTasks)
	router.Get("/list/{id}", h.GetTaskById)
	router.Post("/create", h.CreateTask)
	router.Put("/update/{id}", h.UpdateTask)
	router.Delete("/delete/{id}", h.DeleteTask)
}

func (h *TodoHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.LoggerFromContext(ctx)
	tasks, err := h.service.GetAllTasks(ctx)
	if err != nil {
		logger.Error("failed to get tasks", "error", err)
		response.Json(w, map[string]string{"error": "task not found"}, http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, tasks)
}

func (h *TodoHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.LoggerFromContext(ctx)
	logger.With("method", "get by id")
	id := chi.URLParam(r, "id")

	task, err := h.service.GetTaskById(ctx, id)
	if err != nil {
		logger.Error("Failed to get task", "task_id", id, "error", err)
		response.Json(w, map[string]string{"error": "task not found"}, http.StatusNotFound)
		return
	}

	response.Json(w, task, http.StatusOK)
}
func (h *TodoHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.LoggerFromContext(ctx)
	var requestData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.Error("Invalid request format", "error", err)
		response.Json(w, map[string]string{"error": "invalid request format"}, http.StatusBadRequest)
		return
	}

	// Вызываем gRPC клиент
	task, err := h.service.CreateTask(ctx, requestData.Title, requestData.Description)
	if err != nil {
		logger.Error("Failed to create task", "error", err)
		response.Json(w, map[string]string{"error": "failed to create task"}, http.StatusInternalServerError)
		return
	}

	// Отправляем событие в Kafka
	if h.producer != nil {
		msg := struct {
			Event string      `json:"event"`
			Task  interface{} `json:"task"`
		}{
			Event: "task_created",
			Task:  task,
		}
		if data, err := json.Marshal(msg); err == nil {
			h.producer.SendEvent(ctx, string(data))
		}
	}

	response.Json(w, task, http.StatusCreated)
}

func (h *TodoHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := context.LoggerFromContext(ctx)

	id := chi.URLParam(r, "id")

	var requestData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		logger.Error("Invalid request format", "error", err)
		response.Json(w, map[string]string{"error": "invalid request format"}, http.StatusBadRequest)
		return
	}

	task, err := h.service.UpdateTask(ctx, id, requestData.Title, requestData.Description, requestData.Status)
	if err != nil {
		logger.Error("failed to update task", "id", id, "error", err)
		response.Json(w, map[string]string{"error": "failed to update task"}, http.StatusInternalServerError)

		return
	}
	response.Json(w, task, http.StatusOK)
}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := context.LoggerFromContext(ctx)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid task ID", "id", idStr, "error", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid task ID"})
		return
	}

	err = h.service.DeleteTask(ctx, int32(id))
	if err != nil {
		log.Error("failed to delete task", "id", id, "error", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to delete task"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "task deleted successfully",
	})
}
