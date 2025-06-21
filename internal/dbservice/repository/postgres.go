package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/SteepTaq/todo_project/internal/dbservice/domain"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

type PostgresConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	MaxConns    int
	MaxIdleTime time.Duration
}

func NewPostgresRepo(dsn string, maxConns int, maxIdleTime time.Duration, logger *slog.Logger) (*PostgresRepo, error) {

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	poolCfg.MaxConns = int32(maxConns)
	poolCfg.MaxConnIdleTime = maxIdleTime
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepo{
		pool: pool,
		log:  logger.With("component", "postgres_repo"),
	}, nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}
func (r *PostgresRepo) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at 
              FROM tasks WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var task domain.Task
	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *PostgresRepo) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `INSERT INTO tasks (id, title, description, status, created_at) 
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id, title, description, status, created_at, updated_at`

	row := r.pool.QueryRow(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
	)

	var createdTask domain.Task
	if err := row.Scan(
		&createdTask.ID,
		&createdTask.Title,
		&createdTask.Description,
		&createdTask.Status,
		&createdTask.CreatedAt,
		&createdTask.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return &createdTask, nil
}
func (r *PostgresRepo) UpdateTask(ctx context.Context, todo *domain.Task) error {
	_, err := r.pool.Exec(ctx, "UPDATE todos SET title = $1, description = $2, status = $3, updated_at = $4 WHERE id = $5", todo.Title, todo.Description, todo.Status, todo.UpdatedAt, todo.ID)
	return err
}

func (r *PostgresRepo) DeleteTask(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM todos WHERE id = $1", id)
	return err
}

func (r *PostgresRepo) GetAllTasks(ctx context.Context) ([]*domain.Task, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*domain.Task
	for rows.Next() {
		var todo domain.Task
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}
