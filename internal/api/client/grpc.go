package client

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/SteepTaq/todo_project/internal/api/domain"
	pb "github.com/SteepTaq/todo_project/pkg/proto/gen/todo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// DBClient представляет клиент для взаимодействия с gRPC сервисом
type DBClient struct {
	conn    *grpc.ClientConn
	timeout time.Duration
	client  pb.TodoServiceClient
	logger  *slog.Logger
}

// NewDBClient создает новый экземпляр DBClient
func NewDBClient(target string, timeout time.Duration, logger *slog.Logger) (*DBClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("gRPC connection failed", "error", err, "target", target)
		return nil, err
	}
	return &DBClient{
		conn:    conn,
		client:  pb.NewTodoServiceClient(conn),
		timeout: timeout,
		logger:  logger,
	}, nil
}

// Закрыть соединение
func (c *DBClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	c.logger.Info("gRPC connection closed")

}

func (c *DBClient) GetAllTasks(ctx context.Context) (interface{}, error) {
	// TODO: Реализовать вызов gRPC
	return []interface{}{}, nil
}

func (c *DBClient) GetTaskById(ctx context.Context, id string) (*domain.Task, error) {
	const method = "GetTaskById"
	start := time.Now()
	c.logger.DebugContext(ctx, "gRPC call started",
		"method", method,
		"task_id", id,
	)
	req := &pb.GetTaskRequest{Id: id}

	_, err := c.client.GetTask(ctx, req)
	if err != nil {
		grpcErr := handleGRPCError(err)
		c.logger.ErrorContext(ctx, "gRPC call failed",
			"method", method,
			"task_id", id,
			"error", grpcErr,
			"duration", time.Since(start),
		)
		return nil, grpcErr
	}

	c.logger.DebugContext(ctx, "gRPC call completed",
		"method", method,
		"task_id", id,
		"duration", time.Since(start),
	)

	return nil, nil
	// return &domain.Task{
	// 	ID:          resp.Task.GetTaskId(),
	// 	Title:       resp.Task.GetTitle(),
	// 	Description: resp.Task.GetDescription(),
	// 	Status:      pb.TaskStatus_name[int32(resp.Task.Status)],
	// 	CreatedAt:   resp.Task.GetCreatedAt().AsTime(),
	// 	UpdatedAt:   resp.Task.GetUpdatedAt().AsTime(),
	// }, nil
}

func (c *DBClient) CreateTask(ctx context.Context, title, description string) (*domain.Task, error) {
	start := time.Now()
	const method = "CreateTask"
	c.logger.DebugContext(ctx, "gRPC call started",
		"method", method, "title", title)

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.CreateTaskRequest{
		Task: &pb.Task{
			Title:       title,
			Description: description,
			Status:      pb.TaskStatus_TASK_STATUS_PENDING,
		},
	}

	resp, err := c.client.CreateTask(ctx, req)
	if err != nil {
		c.logger.ErrorContext(ctx, "gRPC call failed", "error", err, "duration", time.Since(start))
		return nil, handleGRPCError(err)
	}

	// Проверяем успешность операции
	if !resp.Success {
		c.logger.Warn("DB service returned unsuccessful response", "duration", time.Since(start))
		return nil, errors.New("failed to create task")
	}
	// Преобразуем protobuf ответ в доменную модель
	task := &domain.Task{
		ID:          resp.Task.TaskId,
		Title:       resp.Task.Title,
		Description: resp.Task.Description,
		Status:      resp.Task.Status.String(),
		CreatedAt:   resp.Task.CreatedAt.AsTime(),
	}

	if resp.Task.UpdatedAt != nil {
		task.UpdatedAt = resp.Task.UpdatedAt.AsTime()
	}
	c.logger.DebugContext(ctx, "Task created",
		"method", method, "task_id", task.ID, "duration", time.Since(start))

	return task, nil
}

func (c *DBClient) UpdateTask(ctx context.Context, id int32, title, description string, completed bool) error {
	// TODO: Реализовать вызов gRPC
	return nil
}

func (c *DBClient) DeleteTask(ctx context.Context, id int32) error {
	// TODO: Реализовать вызов gRPC
	return nil
}

func handleGRPCError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return domain.ErrTaskNotFound
	case codes.AlreadyExists:
		return domain.ErrTaskAlreadyExists
	case codes.InvalidArgument:
		return domain.ErrInvalidInput
	case codes.DeadlineExceeded:
		return domain.ErrRequestTimeout
	case codes.Unavailable:
		return domain.ErrServiceUnavailable
	default:
		return err
	}
}
