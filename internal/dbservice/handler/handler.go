package handler

import (
	"context"

	"github.com/SteepTaq/todo_project/internal/dbservice/repository"
	todo "github.com/SteepTaq/todo_project/pkg/proto/gen_proto/todo"
)

type TodoService interface {
	GetTask(ctx context.Context, req *todo.GetTaskRequest) (*todo.GetTaskResponse, error)
	CreateTask(ctx context.Context, req *todo.CreateTaskRequest) (*todo.CreateTaskResponse, error)
	UpdateTask(ctx context.Context, req *todo.UpdateTaskRequest) (*todo.UpdateTaskResponse, error)
	DeleteTask(ctx context.Context, req *todo.DeleteTaskRequest) (*todo.DeleteTaskResponse, error)
	GetAllTasks(ctx context.Context, req *todo.GetAllTasksRequest) (*todo.GetAllTasksResponse, error)
}

type TodoService struct {
	repo repository.PostgresRepository

}

func GRPCGetTask(ctx context.Context, req *todo.GetTaskRequest) (*todo.GetTaskResponse, error) {

}

func GRPCCreateTask(ctx context.Context, req *todo.CreateTaskRequest) (*todo.CreateTaskResponse, error) {

}

func GRPCUpdateTask(ctx context.Context, req *todo.UpdateTaskRequest) (*todo.UpdateTaskResponse, error) {

}
