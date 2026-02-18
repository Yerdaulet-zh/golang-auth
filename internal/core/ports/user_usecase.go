package ports

import (
	"context"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string) error
	VerifyUserEmail(ctx context.Context, token string) error
	ResendEmailVerificationToken(ctx context.Context, email string) error
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}
