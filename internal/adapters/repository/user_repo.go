package repository

import (
	"context"
	"errors"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *gorm.DB
	logger ports.Logger
}

func NewUserRepository(db *gorm.DB, logger ports.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*repouser.User, error) {
	userRecord, err := gorm.G[repouser.User](repo.db).Where("email = ?", email).Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		repo.logger.Error("Repository", "Error while querying user by email", "error", err)
		return nil, domain.ErrDatabaseInternalError
	}
	return &userRecord, nil
}

func (repo *UserRepository) CreateUserWithCredentials(ctx context.Context, req ports.UserAndCredentialsRequest) error {
	repoUser := repouser.User{
		Email:        req.Email,
		UserStatus:   req.UserStatus,
		IsMFAEnabled: req.IsMFAEnabled,
	}
	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := gorm.G[repouser.User](repo.db).Create(ctx, &repoUser); err != nil {
			repo.logger.Error("Repository", "Error while insert into user table", "error", err)
			return domain.ErrDatabaseInternalError
		}

		repoCred := repouser.UserCredentials{
			UserID:       repoUser.ID,
			PasswordHash: req.PasswordHash,
		}
		if err := gorm.G[repouser.UserCredentials](repo.db).Create(ctx, &repoCred); err != nil {
			repo.logger.Error("Repository", "Error while insert into user_credentials table", "error", err)
			return domain.ErrDatabaseInternalError
		}
		return nil
	}); err != nil {
		repo.logger.Error("Repository", "Error from transaction | CreateUSerWithCredentials", "error", err)
		return domain.ErrDatabaseInternalError
	}
	return nil
}
