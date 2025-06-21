package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SteepTaq/todo_project/internal/dbservice/config"
	"github.com/SteepTaq/todo_project/internal/dbservice/repository"
	"github.com/SteepTaq/todo_project/internal/dbservice/server"
	"github.com/SteepTaq/todo_project/internal/dbservice/service"
	"github.com/SteepTaq/todo_project/pkg/logger"
	todov1 "github.com/SteepTaq/todo_project/pkg/proto/gen/todo"
	"google.golang.org/grpc"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация логгера
	log := logger.Setup(cfg.Logger.Level)
	slog.SetDefault(log)

	// Создание контекста для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск приложения
	if err := app(ctx, cfg, log); err != nil {
		log.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func app(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	// Формируем DSN для Postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	// Инициализация репозиториев
	pgRepo, err := repository.NewPostgresRepo(
		dsn,
		cfg.Postgres.MaxConns,
		cfg.Postgres.MaxIdleTime,
		log,
	)
	if err != nil {
		return fmt.Errorf("failed to create Postgres repo: %w", err)
	}
	defer pgRepo.Close()

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisRepo, err := repository.NewRedisRepo(
		redisAddr,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.CacheTTL,
		log,
	)
	if err != nil {
		return fmt.Errorf("failed to create Redis repo: %w", err)
	}
	defer redisRepo.Close()

	// Создание сервиса с кеширующим слоем
	taskService := service.NewTaskService(pgRepo, redisRepo, log)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer(
		grpc.ConnectionTimeout(cfg.GRPC.Timeout),
	)

	// Регистрация сервиса
	todov1.RegisterTodoServiceServer(grpcServer, server.NewGRPCServer(taskService))

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", cfg.GRPC.Target)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	go func() {
		log.Info("starting gRPC server", "address", cfg.GRPC.Target)
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server failed", "error", err)
		}
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Info("shutting down server")

	// Graceful shutdown
	grpcServer.GracefulStop()
	log.Info("server stopped gracefully")

	return nil
}
