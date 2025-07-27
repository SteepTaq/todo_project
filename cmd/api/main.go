package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SteepTaq/todo_project/internal/api/client"
	"github.com/SteepTaq/todo_project/internal/api/config"
	"github.com/SteepTaq/todo_project/internal/api/handler"
	"github.com/SteepTaq/todo_project/internal/api/kafka"
	ctxLog "github.com/SteepTaq/todo_project/pkg/context"
	"github.com/SteepTaq/todo_project/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация логгера
	log := logger.Setup(cfg.Logger.Level)
	slog.SetDefault(log)
	log.Info("Configuration loaded",
		"http_port", cfg.HTTP.Port,
		"grpc_target", cfg.GRPC.Target,
		"kafka_topic", cfg.Kafka.Brokers,
		"logger_level", cfg.Logger.Level,
	)
	// Создание контекста для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск приложения
	if err := App(ctx, cfg, log); err != nil {
		log.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func App(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	// Инициализация gRPC клиента (заглушка, реализация в client/grpc.go)
	dbClient, err := client.NewDBClient(cfg.GRPC.Target, cfg.GRPC.Timeout, log)
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}
	defer dbClient.Close()

	// Создание HTTP роутера
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	//r.Use(middleware.GetStructuredLogger(log))
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestLogger := log.With(
				"request_id", middleware.GetReqID(r.Context()),
				"method", r.Method,
				"path", r.URL.Path,
			)
			ctx := ctxLog.WithLogger(r.Context(), requestLogger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// Создать Kafka-продюсер
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer producer.Close()

	// Инициализация и регистрация обработчиков
	todoHandler := handler.NewTodoHandler(cfg, dbClient, producer)
	todoHandler.RegisterRoutes(r)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		log.Info("ok")
	})

	// HTTP сервер
	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Запуск сервера в горутине
	serverErr := make(chan error, 1)
	go func() {
		log.Info("starting HTTP server", "port", cfg.HTTP.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		log.Info("shutting down server")
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Info("server stopped gracefully")
	return nil
}
