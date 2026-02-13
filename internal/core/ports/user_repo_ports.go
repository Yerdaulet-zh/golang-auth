package ports

import (
	"context"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
)

type UserRepoPorts interface {
	GetUserByEmail(ctx context.Context, email string) (*repouser.User, error)
	CreateUserWithCredentials(ctx context.Context, req UserAndCredentialsRequest) error
}
