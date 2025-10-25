package user

import "time"

// Event types for event-driven communication
const (
	EventTypeUserCreated = "user.created"
	EventTypeUserUpdated = "user.updated"
	EventTypeUserDeleted = "user.deleted"
)

// UserCreatedEvent is published when a new user is created
type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUpdatedEvent is published when a user is updated
type UserUpdatedEvent struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserDeletedEvent is published when a user is deleted
type UserDeletedEvent struct {
	UserID    string    `json:"user_id"`
	DeletedAt time.Time `json:"deleted_at"`
}
