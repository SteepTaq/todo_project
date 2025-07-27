package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9094" // fallback для локального запуска
	}

	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "events" // fallback
	}

	logFile := os.Getenv("LOG_FILE")
	if logFile == "" {
		logFile = "events.log" // fallback
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokers},
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	logger.Printf("Worker started, connecting to %s, topic: %s", brokers, topic)
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				logger.Println("Shutting down worker...")
				return
			}
			logger.Printf("read error: %v", err)
			continue
		}
		logger.Printf("Received: %s", string(m.Value))
	}
}
