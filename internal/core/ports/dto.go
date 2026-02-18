package ports

import (
	"time"

	"github.com/google/uuid"
)

type UserAndCredentialsRequest struct {
	Email                  string
	UserStatus             string
	IsMFAEnabled           bool
	PasswordHash           string
	EmailVerificationToken string
	TokenExpiration        time.Time
}

type CreateUserSessionRequest struct {
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Device    string
	Token     string
	ExpiresAt time.Time
}

type CreateUserSessionResponse struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Token      string
	ExpiresAt  time.Time
	LastActive time.Time
}

type LoginRequest struct {
	Email     string
	Password  string
	IPAddress string
	UserAgent string
	Device    string
}

type LoginResponse struct {
	SessionID    uuid.UUID
	UserID       uuid.UUID
	RefreshToken string
	AccessToken  string

	RefreshTokenExpiresAt time.Time
	AccessTokenExpiresAt  time.Time
}
