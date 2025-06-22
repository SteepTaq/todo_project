package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/SteepTaq/todo_project/internal/dbservice/config"
	"github.com/SteepTaq/todo_project/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Парсинг аргументов командной строки
	action := flag.String("action", "up", "Migration action (up, down, force)")
	steps := flag.Int("steps", 0, "Number of steps for migration")
	version := flag.String("version", "", "Version for force migration")
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация логгера
	log := logger.Setup(cfg.Logger.Level)
	slog.SetDefault(log)

	// Формирование DSN
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)

	// Подключение к БД для проверки
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Error("Database ping failed", "error", err)
		os.Exit(1)
	}

	// Инициализация мигратора
	m, err := migrate.New(
		"file://internal/dbservice/migrations",
		dsn,
	)
	if err != nil {
		log.Error("Failed to initialize migrator", "error", err)
		os.Exit(1)
	}

	// Выполнение действия
	switch *action {
	case "up":
		if *steps > 0 {
			log.Info("Applying migrations", "steps", *steps)
			err = m.Steps(*steps)
		} else {
			log.Info("Applying all migrations")
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			log.Info("Reverting migrations", "steps", *steps)
			err = m.Steps(-*steps)
		} else {
			log.Info("Reverting last migration")
			err = m.Down()
		}
	case "force":
		if *version == "" {
			log.Error("Version is required for force migration")
			os.Exit(1)
		}
		v, errConv := strconv.Atoi(*version)
		if errConv != nil {
			log.Error("Invalid version number", "version", *version)
			os.Exit(1)
		}
		log.Info("Forcing migration version", "version", v)
		err = m.Force(v)
	default:
		log.Error("Unknown action", "action", *action)
		os.Exit(1)
	}

	// Обработка результата
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Info("No migrations applied - database is up to date")
		} else {
			log.Error("Migration failed", "error", err)
			os.Exit(1)
		}
	}

	log.Info("Migration completed successfully")
}
