package ports

import (
	"context"
)

type UserUseCase interface {
	Register(ctx context.Context, email, password string) error
}
