package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			
			Balancer: &kafka.LeastBytes{},
		},
		topic: topic,
	}
}

func (p *Producer) SendEvent(ctx context.Context, message string) error {
	log.Printf("Attempting to send message to Kafka: %s", message)

	msg := kafka.Message{
		Value: []byte(message),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("failed to send kafka message: %v", err)
	} else {
		log.Printf("Message sent successfully to Kafka, topic: %s (using LeastBytes balancer)", p.topic)
	}
	return err
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}
