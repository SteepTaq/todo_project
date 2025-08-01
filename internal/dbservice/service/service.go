package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/SteepTaq/todo_project/internal/dbservice/domain"
	"github.com/google/uuid"
)

type TaskService struct {
	storage TaskRepository
	cache   TaskCache
	log     *slog.Logger
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	GetAllTasks(ctx context.Context) ([]*domain.Task, error)
	UpdateTask(ctx context.Context, tasks *domain.Task) (*domain.Task, error)
}

type TaskCache interface {
	SetTask(ctx context.Context, task *domain.Task) error
	GetTask(ctx context.Context, id string) (*domain.Task, error)
}

func NewTaskService(storage TaskRepository, cache TaskCache, logger *slog.Logger) *TaskService {
	return &TaskService{
		storage: storage,
		cache:   cache,
		log:     logger.With("component", "task_service"),
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	start := time.Now()

	if task.Title == "" {
		return nil, domain.ErrInvalidInput
	}
	newID := uuid.New().String()

	newTask := &domain.Task{
		ID:          newID,
		Title:       task.Title,
		Description: task.Description,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTask, err := s.storage.CreateTask(ctx, newTask)
	if err != nil {
		s.log.Error("failed to create task", "error", err)
		return nil, err
	}

	if err := s.cache.SetTask(ctx, createdTask); err != nil {
		s.log.Warn("failed to cache task", "task_id", createdTask.ID, "error", err)
	}

	s.log.Info("task created",
		"task_id", createdTask.ID,
		"duration", time.Since(start))

	return createdTask, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	start := time.Now()

	if task.Title == "" {
		return nil, domain.ErrInvalidInput
	}
	
	newTask := &domain.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedTask, err := s.storage.UpdateTask(ctx, newTask)
	if err != nil {
		s.log.Error("failed to create task", "error", err)
		return nil, err
	}

	if err := s.cache.SetTask(ctx, updatedTask); err != nil {
		s.log.Warn("failed to cache task", "task_id", updatedTask.ID, "error", err)
	}

	s.log.Info("task updated",
		"task_id", updatedTask.ID,
		"duration", time.Since(start))

	return updatedTask, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	start := time.Now()

	if cachedTask, err := s.cache.GetTask(ctx, id); err == nil {
		s.log.Debug("task retrieved from cache",
			"task_id", id,
			"duration", time.Since(start))
		return cachedTask, nil
	}

	task, err := s.storage.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			s.log.Warn("task not found", "task_id", id)
		} else {
			s.log.Error("failed to get task", "task_id", id, "error", err)
		}
		return nil, err
	}

	if err := s.cache.SetTask(ctx, task); err != nil {
		s.log.Warn("failed to cache task", "task_id", id, "error", err)
	}

	s.log.Debug("task retrieved from storage",
		"task_id", id,
		"duration", time.Since(start))

	return task, nil
}
func (s *TaskService) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	start := time.Now()

	tasks, err := s.storage.GetAllTasks(ctx)
	if err != nil {
		s.log.Error("failed to get tasks", "error", err)
		return nil, err
	}

	s.log.Debug("task retrieved from storage",
		"duration", time.Since(start))

	return tasks, nil
}
