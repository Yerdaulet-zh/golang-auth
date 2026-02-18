package ports

import (
	"context"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string) error
	VerifyUserEmail(ctx context.Context, token string) error
	ResendEmailVerificationToken(ctx context.Context, email string) error
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, session_id uuid.UUID) error
	DeleteAccount(ctx context.Context, user_id uuid.UUID) error
}
