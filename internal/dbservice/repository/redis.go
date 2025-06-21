package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/SteepTaq/todo_project/internal/dbservice/domain"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
	log    *slog.Logger
	ttl    time.Duration
}

func NewRedisRepo(addr, password string, db int, ttl time.Duration, logger *slog.Logger) (*RedisRepo, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisRepo{
		client: client,
		log:    logger.With("component", "redis_repo"),
		ttl:    ttl,
	}, nil
}

func (r *RedisRepo) Close() {
	r.client.Close()
}

func (r *RedisRepo) SetTask(ctx context.Context, task *domain.Task) error {
	key := "task:" + task.ID
	value, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	if err := r.client.Set(ctx, key, value, r.ttl).Err(); err != nil {
		return fmt.Errorf("failed to set task in Redis: %w", err)
	}

	return nil
}

func (r *RedisRepo) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	key := "task:" + id
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("failed to get task from Redis: %w", err)
	}

	var task domain.Task
	if err := json.Unmarshal([]byte(value), &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}
