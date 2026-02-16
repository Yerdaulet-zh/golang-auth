package ports

import (
	"context"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string) error
	VerifyUserEmail(ctx context.Context, token string) error
}
