package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        string
	Email     string
	Name      string
	Password  string // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new user with hashed password
func NewUser(email, name, password string) (*User, error) {
	if email == "" {
		return nil, ErrInvalidEmail
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		Email:     email,
		Name:      name,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// CheckPassword verifies if the provided password is correct
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)
}

// UpdateProfile updates user profile fields
func (u *User) UpdateProfile(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	u.Name = name
	u.UpdatedAt = time.Now()
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	return string(bytes), err
}
