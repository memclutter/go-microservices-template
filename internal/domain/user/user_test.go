package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		password string
		wantErr  error
	}{
		{
			name:     "valid user",
			email:    "test@example.com",
			userName: "Test User",
			password: "password123",
			wantErr:  nil,
		},
		{
			name:     "invalid email",
			email:    "",
			userName: "Test User",
			password: "password123",
			wantErr:  ErrInvalidEmail,
		},
		{
			name:     "invalid name",
			email:    "test@example.com",
			userName: "",
			password: "password123",
			wantErr:  ErrInvalidName,
		},
		{
			name:     "weak password",
			email:    "test@example.com",
			userName: "Test User",
			password: "pass",
			wantErr:  ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.userName, tt.password)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, user)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tt.email, user.Email)
			assert.Equal(t, tt.userName, user.Name)
			assert.NotEmpty(t, user.Password)
			assert.NotEqual(t, tt.password, user.Password) // Password должен быть хэширован
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user, err := NewUser("test@example.com", "Test", "password123")
	require.NoError(t, err)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "correct password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "incorrect password",
			password: "wrongpassword",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.CheckPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_UpdateProfile(t *testing.T) {
	user, err := NewUser("test@example.com", "Test", "password123")
	require.NoError(t, err)

	oldUpdatedAt := user.UpdatedAt

	err = user.UpdateProfile("New Name")
	require.NoError(t, err)
	assert.Equal(t, "New Name", user.Name)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))

	err = user.UpdateProfile("")
	assert.ErrorIs(t, err, ErrInvalidName)
}
