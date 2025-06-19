package main

import (
	"database/sql"

	"github.com/IBM/sarama"
	"github.com/SteepTaq/todo_project/internal/config"
	"github.com/SteepTaq/todo_project/internal/dbservice/client"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type App struct {
	db     *sql.DB
	redis  *redis.Client
	router *chi.Mux
}

func NewApp(kafkaProducer *sarama.Producer, dbClient *client.DBServiceClient, redisClient *redis.Client) *App {
	return &App{
		db:     dbClient,
		redis:  redisClient,
		router: chi.NewRouter(),
	}
}

func main() {
	cfg := config.Load()
	kafkaProducer := kafka.NewProducer(cfg.Kafka)
	dbClient := client.NewDBServiceClient(cfg.DBService)
	app := NewApp(kafkaProducer, dbClient)
	app.Run(cfg.HTTP.Port)
}
