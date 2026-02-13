package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	UserStatus   string
	IsMFAEnabled bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// If DeletedAt.IsZero() == true, the user is not deleted.
	DeletedAt   time.Time
	Credentials UserCredentials
}

type UserCredentials struct {
	UserID             uuid.UUID
	PasswordHash       string
	LastPasswordChange time.Time
	UpdatedAt          time.Time
}
