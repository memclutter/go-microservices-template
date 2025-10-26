package rabbitmq

import (
	"context"
	"fmt"

	"github.com/memclutter/go-microservices-template/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EventHandler is a function that handles incoming events
type EventHandler func(ctx context.Context, eventType string, payload []byte) error

// Consumer consumes messages from RabbitMQ
type Consumer struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	queue    string
	handlers map[string]EventHandler
	logger   *logger.Logger
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(url, queueName string, routingKeys []string, log *logger.Logger) (*Consumer, error) {
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

	// Declare exchange (idempotent)
	err = ch.ExchangeDeclare(
		EventsExchange,
		ExchangeTypeTopic,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange with routing keys
	for _, routingKey := range routingKeys {
		err = ch.QueueBind(
			q.Name,         // queue name
			routingKey,     // routing key
			EventsExchange, // exchange
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	log.WithFields(map[string]any{
		"queue":        queueName,
		"routing_keys": routingKeys,
	}).Info("RabbitMQ consumer initialized")

	return &Consumer{
		conn:     conn,
		ch:       ch,
		queue:    q.Name,
		handlers: make(map[string]EventHandler),
		logger:   log,
	}, nil
}

// RegisterHandler registers an event handler for a specific event type
func (c *Consumer) RegisterHandler(eventType string, handler EventHandler) {
	c.handlers[eventType] = handler
	c.logger.WithField("event_type", eventType).Debug("Event handler registered")
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	// Set QoS (prefetch count)
	err := c.ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	msgs, err := c.ch.Consume(
		c.queue, // queue
		"",      // consumer tag
		false,   // auto-ack (manual ack for reliability)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.logger.Info("Started consuming messages")

	// Process messages
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer stopped by context")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Warn("Message channel closed")
				return fmt.Errorf("message channel closed")
			}
			c.handleMessage(ctx, msg)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	eventType := msg.RoutingKey

	c.logger.WithFields(map[string]any{
		"event_type": eventType,
		"size":       len(msg.Body),
	}).Debug("Received message")

	// Find handler
	handler, ok := c.handlers[eventType]
	if !ok {
		c.logger.WithField("event_type", eventType).Warn("No handler registered for event type")
		msg.Nack(false, false) // Reject message
		return
	}

	// Handle event
	err := handler(ctx, eventType, msg.Body)
	if err != nil {
		c.logger.WithError(err).WithField("event_type", eventType).Error("Failed to handle event")
		msg.Nack(false, true) // Requeue message
		return
	}

	// Acknowledge message
	if err := msg.Ack(false); err != nil {
		c.logger.WithError(err).Error("Failed to acknowledge message")
	}
}

// Close closes the RabbitMQ connection
func (c *Consumer) Close() error {
	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			c.logger.WithError(err).Warn("Failed to close RabbitMQ channel")
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			c.logger.WithError(err).Warn("Failed to close RabbitMQ connection")
			return err
		}
	}
	c.logger.Info("RabbitMQ consumer closed")
	return nil
}
