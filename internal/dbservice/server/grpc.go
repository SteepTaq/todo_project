package server

import (
	"context"
	"errors"

	"github.com/SteepTaq/todo_project/internal/dbservice/domain"
	"github.com/SteepTaq/todo_project/internal/dbservice/service"
	todov1 "github.com/SteepTaq/todo_project/pkg/proto/gen/todo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	todov1.UnimplementedTodoServiceServer
	service *service.TaskService
}

func NewGRPCServer(service *service.TaskService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) CreateTask(ctx context.Context, req *todov1.CreateTaskRequest) (*todov1.CreateTaskResponse, error) {
	domainTask := &domain.Task{
		Title:       req.Task.GetTitle(),
		Description: req.Task.GetDescription(),
		Status:      req.Task.GetStatus().String(),
	}

	newTask, err := s.service.CreateTask(ctx, domainTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbTask := &todov1.Task{
		TaskId:      newTask.ID,
		Title:       newTask.Title,
		Description: newTask.Description,
		Status:      todov1.TaskStatus(todov1.TaskStatus_value[newTask.Status]),
		CreatedAt:   timestamppb.New(newTask.CreatedAt),
	}

	if !newTask.UpdatedAt.IsZero() {
		pbTask.UpdatedAt = timestamppb.New(newTask.UpdatedAt)
	}

	return &todov1.CreateTaskResponse{
		Success: true,
		Task:    pbTask,
	}, nil
}

func (s *GRPCServer) GetTask(ctx context.Context, req *todov1.GetTaskRequest) (*todov1.GetTaskResponse, error) {
	task, err := s.service.GetTask(ctx, req.GetTaskId())
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	todov1Task := &todov1.Task{
		TaskId:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      todov1.TaskStatus(todov1.TaskStatus_value[task.Status]),
		CreatedAt:   timestamppb.New(task.CreatedAt),
	}

	if !task.UpdatedAt.IsZero() {
		todov1Task.UpdatedAt = timestamppb.New(task.UpdatedAt)
	}

	return &todov1.GetTaskResponse{
		Task: todov1Task,
	}, nil
}
