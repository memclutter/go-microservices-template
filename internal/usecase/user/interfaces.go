package user

import "context"

// EventPublisher defines interface for publishing domain events
type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload interface{}) error
}
