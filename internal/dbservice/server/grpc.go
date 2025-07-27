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
		UpdatedAt:   timestamppb.New(newTask.UpdatedAt),
	}

	return &todov1.CreateTaskResponse{
		Success: true,
		Task:    pbTask,
	}, nil
}

func (s *GRPCServer) GetTask(ctx context.Context, req *todov1.GetTaskRequest) (*todov1.GetTaskResponse, error) {
	newTask, err := s.service.GetTask(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return nil, status.Error(codes.NotFound, "task not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbTask := &todov1.Task{
		TaskId:      newTask.ID,
		Title:       newTask.Title,
		Description: newTask.Description,
		Status:      todov1.TaskStatus(todov1.TaskStatus_value[newTask.Status]),
		CreatedAt:   timestamppb.New(newTask.CreatedAt),
		UpdatedAt:   timestamppb.New(newTask.UpdatedAt),
	}

	return &todov1.GetTaskResponse{
		Task: pbTask,
	}, nil
}
func (s *GRPCServer) GetAllTasks(ctx context.Context, req *todov1.GetAllTasksRequest) (*todov1.GetAllTasksResponse, error) {
	newTasks, err := s.service.GetAllTasks(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pbTasks := make([]*todov1.Task, 0, len(newTasks))
	for _, task := range newTasks {
		t := todov1.Task{
			TaskId:      task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      todov1.TaskStatus(todov1.TaskStatus_value[task.Status]),
			CreatedAt:   timestamppb.New(task.CreatedAt),
			UpdatedAt:   timestamppb.New(task.UpdatedAt),
		}
		pbTasks = append(pbTasks, &t)
	}

	return &todov1.GetAllTasksResponse{
		Tasks: pbTasks,
	}, nil
}

func (s *GRPCServer) UpdateTask(ctx context.Context, req *todov1.UpdateTaskRequest) (*todov1.UpdateTaskResponse, error) {
	domainTask := &domain.Task{
		ID:          req.Task.GetTaskId(),
		Title:       req.Task.GetTitle(),
		Description: req.Task.GetDescription(),
	}
	switch req.Task.GetStatus() {
	case todov1.TaskStatus_TASK_STATUS_PENDING:
		domainTask.Status = "pending"
	case todov1.TaskStatus_TASK_STATUS_IN_PROGRESS:
		domainTask.Status = "in_progress"
	case todov1.TaskStatus_TASK_STATUS_COMPLETED:
		domainTask.Status = "completed"
	}
	newTask, err := s.service.UpdateTask(ctx, domainTask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbTask := &todov1.Task{
		TaskId:      newTask.ID,
		Title:       newTask.Title,
		Description: newTask.Description,
		CreatedAt:   timestamppb.New(newTask.CreatedAt),
		UpdatedAt:   timestamppb.New(newTask.UpdatedAt),
	}
	switch newTask.Status {
	case "pending":
		pbTask.Status = todov1.TaskStatus(0)
	case "in_progress":
		pbTask.Status = todov1.TaskStatus(1)
	case "completed":
		pbTask.Status = todov1.TaskStatus(2)
	}
	return &todov1.UpdateTaskResponse{
		Task: pbTask,
	}, nil
}
