package repository

import (
	"context"
	"errors"
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	userverification "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_verification"
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
		if err := gorm.G[repouser.User](tx).Create(ctx, &repoUser); err != nil {
			repo.logger.Error("Repository", "Error while insert into user table", "error", err)
			return domain.ErrDatabaseInternalError
		}

		repoCred := repouser.UserCredentials{
			UserID:       repoUser.ID,
			PasswordHash: req.PasswordHash,
		}
		if err := gorm.G[repouser.UserCredentials](tx).Create(ctx, &repoCred); err != nil {
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

var (
	ErrTokenNotFound     = errors.New("User verification token is not found, invalid token")
	ErrTokenExpired      = errors.New("Token is expired")
	ErrInvalidTokenState = errors.New("Invalid token")
)

func (repo *UserRepository) VerifyUserEmail(ctx context.Context, token string) error {
	record, err := gorm.G[userverification.UserVerification](repo.db).Where("token = ?", token).Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTokenNotFound
		}
		repo.logger.Error(
			"Repository", "Error while querying user email verification token from the user_verification table", "error", err, "token", token,
		)
		return domain.ErrDatabaseInternalError
	}

	if record.Status == "active" {
		repo.logger.Info("Repository", "Token already verified", "token", token)
		return nil
	}

	if time.Now().After(record.ExpiresAt) {
		return ErrTokenExpired
	}

	if record.Status != "pending" {
		return ErrInvalidTokenState
	}

	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[userverification.UserVerification](tx).Where("id = ?", record.ID).Update(ctx, "status", "active")
		if err != nil {
			repo.logger.Error("Repository", "Failed to update verification status to active", "error", err, "token", token)
			return domain.ErrDatabaseInternalError
		}
		_, err = gorm.G[repouser.User](tx).Where("id = ?", record.UserID).Update(ctx, "user_status", "active")
		if err != nil {
			repo.logger.Error("Repository", "Failed to update user status from user table to active", "error", err, "user_id", record.UserID)
			return domain.ErrDomainInternalError
		}
		return nil
	})
}
