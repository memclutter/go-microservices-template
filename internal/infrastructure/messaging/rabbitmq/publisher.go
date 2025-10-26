package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/memclutter/go-microservices-template/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchange names
	EventsExchange = "microservices.events"

	// Exchange type
	ExchangeTypeTopic = "topic"
)

// Publisher publishes messages to RabbitMQ
type Publisher struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *logger.Logger
}

// NewPublisher creates a new RabbitMQ publisher
func NewPublisher(url string, log *logger.Logger) (*Publisher, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		EventsExchange,    // name
		ExchangeTypeTopic, // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	log.WithField("exchange", EventsExchange).Info("RabbitMQ publisher initialized")

	return &Publisher{
		conn:   conn,
		ch:     ch,
		logger: log,
	}, nil
}

// Publish publishes an event to RabbitMQ
func (p *Publisher) Publish(ctx context.Context, eventType string, payload interface{}) error {
	// Marshal payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Publish message
	err = p.ch.PublishWithContext(
		ctx,
		EventsExchange, // exchange
		eventType,      // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Persistent delivery
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		p.logger.WithError(err).WithField("event_type", eventType).Error("Failed to publish event")
		return fmt.Errorf("failed to publish event: %w", err)
	}

	p.logger.WithFields(map[string]any{
		"event_type": eventType,
		"size":       len(body),
	}).Debug("Event published successfully")

	return nil
}

// Close closes the RabbitMQ connection
func (p *Publisher) Close() error {
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			p.logger.WithError(err).Warn("Failed to close RabbitMQ channel")
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			p.logger.WithError(err).Warn("Failed to close RabbitMQ connection")
			return err
		}
	}
	p.logger.Info("RabbitMQ publisher closed")
	return nil
}
