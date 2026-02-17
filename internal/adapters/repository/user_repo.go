package repository

import (
	"context"
	"errors"
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	userverification "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_verification"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

		// User Token verification
		repoVerify := userverification.UserVerification{
			UserID:    repoUser.ID,
			Token:     req.EmailVerificationToken,
			ExpiresAt: req.TokenExpiration,
		}
		if err := gorm.G[userverification.UserVerification](tx).Create(ctx, &repoVerify); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				repo.logger.Warn("Repository", "Token collision detected, retrying...", "token", repoVerify.Token)
				return domain.ErrTokenCollision
			}

			repo.logger.Error("Repository", "Failed to create verification record", "error", err)
			return domain.ErrDatabaseInternalError
		}
		return nil
	}); err != nil {
		repo.logger.Error("Repository", "Error from transaction | CreateUSerWithCredentials", "error", err)
		return domain.ErrDatabaseInternalError
	}
	return nil
}

func (repo *UserRepository) VerifyUserEmail(ctx context.Context, token string) error {
	record, err := gorm.G[userverification.UserVerification](repo.db).Where("token = ?", token).Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrTokenNotFound
		}
		repo.logger.Error(
			"Repository", "Error while querying user email verification token from the user_verification table", "error", err, "token", token,
		)
		return domain.ErrRepositoryInternalError
	}

	if record.Status == "active" {
		repo.logger.Info("Repository", "Token already verified", "token", token)
		return nil
	}

	if time.Now().After(record.ExpiresAt) {
		return domain.ErrTokenExpired
	}

	if record.Status != "pending" {
		return domain.ErrInvalidTokenState
	}

	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[userverification.UserVerification](tx).Where("id = ?", record.ID).Update(ctx, "status", "consumed")
		if err != nil {
			repo.logger.Error("Repository", "Failed to update verification status to active", "error", err, "token", token)
			return domain.ErrRepositoryInternalError
		}
		_, err = gorm.G[repouser.User](tx).Where("id = ?", record.UserID).Update(ctx, "user_status", "active")
		if err != nil {
			repo.logger.Error("Repository", "Failed to update user status from user table to active", "error", err, "user_id", record.UserID)
			return domain.ErrRepositoryInternalError
		}
		return nil
	})
}

func (repo *UserRepository) GetVerificationByToken(ctx context.Context, token string) (*userverification.UserVerification, error) {
	record, err := gorm.G[userverification.UserVerification](repo.db).Where("token = ?", token).Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTokenNotFound
		}
		return nil, domain.ErrRepositoryInternalError
	}
	return &record, nil
}

func (repo *UserRepository) ConfirmVerification(ctx context.Context, userID uuid.UUID, verificationID uuid.UUID) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update Token to consumed
		if err := tx.Model(&userverification.UserVerification{}).Where("id = ?", verificationID).Update("status", "consumed").Error; err != nil {
			return domain.ErrRepositoryInternalError
		}
		// Update User to active
		if err := tx.Model(&repouser.User{}).Where("id = ?", userID).Update("user_status", "active").Error; err != nil {
			return domain.ErrRepositoryInternalError
		}
		return nil
	})
}

func (repo *UserRepository) GetVerificationByUserID(ctx context.Context, userID uuid.UUID) (*userverification.UserVerification, error) {
	record, err := gorm.G[userverification.UserVerification](repo.db).Where("user_id = ?", userID).Order("created_at DESC").Take(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			repo.logger.Error(domain.LogRepository, "Error for resend verification not found existing token record by user_id", "error", err, "user_id", userID)
			return nil, domain.ErrRepositoryInternalError
		}
		repo.logger.Error(domain.LogRepository, "Unexpected error whike queryng the last token record by user_id", "error", err, "user_id", userID)
		return nil, domain.ErrDatabaseInternalError
	}
	return &record, nil
}

// Bug the old ones are becoming in pedning state: p, i, p, p should be p, i, i, i
func (repo *UserRepository) RotateVerificationToken(ctx context.Context, recordID uuid.UUID, status string, req *userverification.UserVerification) error {
	if err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		_, err := gorm.G[userverification.UserVerification](repo.db).Where("ID = ?", recordID).Update(ctx, "status", status)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				repo.logger.Error(domain.LogRepository, "Error not found user email validation record by id", "error", err, "ID", recordID)
				return domain.ErrRepositoryInternalError
			}
			repo.logger.Error(domain.LogRepository, "Unexpected error while updating the email validation record by id", "error", err, "ID", recordID)
			return domain.ErrDatabaseInternalError
		}

		if err := gorm.G[userverification.UserVerification](tx).Create(ctx, req); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				repo.logger.Warn(domain.LogRepository, "Token collision detected, retrying...", "token", req.Token)
				return domain.ErrTokenCollision
			}

			repo.logger.Error(domain.LogRepository, "Failed to create verification record", "error", err)
			return domain.ErrDatabaseInternalError
		}

		return nil
	}); err != nil {
		return domain.ErrDatabaseInternalError
	}
	return nil
}

func (repo *UserRepository) GetCountsOfVerificationRecordsByUserID(ctx context.Context, user_id uuid.UUID, timeDuration time.Time) (int64, error) {
	count, err := gorm.G[userverification.UserVerification](repo.db).Where("user_id = ? AND created_at >= ?", user_id, timeDuration).Count(ctx, "ID")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			repo.logger.Error(domain.LogRepository, "Error not found user verification records", "error", err, "ID", user_id)
			return 0, domain.ErrRepositoryInternalError
		}
		repo.logger.Error(domain.LogRepository, "Unexpected error not found user verification records", "error", err, "ID", user_id)
		return 0, domain.ErrDatabaseInternalError
	}
	return count, nil
}
