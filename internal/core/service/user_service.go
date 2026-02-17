package service

import (
	"context"
	"errors"
	"strings"
	"time"

	userverification "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user_verification"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSerivce struct {
	repo      ports.UserRepoPorts
	logger    ports.Logger
	publisher ports.EventPublisher
}

func NewUserService(repo ports.UserRepoPorts, logger ports.Logger, publisher ports.EventPublisher) *UserSerivce {
	return &UserSerivce{
		repo:      repo,
		logger:    logger,
		publisher: publisher,
	}
}

func (s *UserSerivce) Register(ctx context.Context, email, password string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	// Initial Validations
	if err := isValidEmail(email); err != nil {
		if errors.Is(err, domain.ErrInvalidEmail) {
			return domain.ErrInvalidEmail
		}
		s.logger.Error("Domain", "Error while checking user email", "error", err)
		return domain.ErrDomainInternalError
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		s.logger.Error("Error while trying to hash password", "error", err)
		return domain.ErrHashingError
	}

	// Check if user exists
	userRecord, err := s.repo.GetUserByEmail(ctx, email)
	if userRecord != nil {
		return domain.ErrUserAlreadyExists
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return domain.ErrDatabaseInternalError
	}

	// Persistence Loop (Retry on Token Collision)
	var finalToken string
	const maxRetries = 3
	committed := false

	for i := 0; i < maxRetries; i++ {
		token, err := GenerateSecureToken()
		if err != nil {
			s.logger.Error(domain.LogService, "Token generation failed", "error", err)
			return domain.ErrDomainInternalError
		}

		repoReq := ports.UserAndCredentialsRequest{
			Email:                  email,
			PasswordHash:           hashedPassword,
			EmailVerificationToken: token,
			TokenExpiration:        time.Now().Add(15 * time.Minute),
		}

		err = s.repo.CreateUserWithCredentials(ctx, repoReq)
		if err == nil {
			finalToken = token
			committed = true
			break
		}

		if !errors.Is(err, domain.ErrTokenCollision) {
			return err
		}

		s.logger.Warn("Service", "Token collision detected, retrying...", "attempt", i+1)
	}

	if !committed {
		s.logger.Error(domain.LogService, "Max retries reached for registration collisions")
		return domain.ErrDatabaseInternalError
	}

	// Publish to Broker (ONLY after successful DB commit)
	if err := s.publisher.PublishUserRegistered(ctx, email, finalToken); err != nil {
		s.logger.Error(domain.LogService, "Broker error after DB commit", "error", err, "email", email)
		// We return an error here, but note that the user is already created in the DB.
		return domain.ErrBrokerInternalError
	}

	return nil
}

func (s *UserSerivce) VerifyUserEmail(ctx context.Context, token string) error {
	// Validation
	record, err := s.repo.GetVerificationByToken(ctx, token)
	if err != nil {
		return domain.ErrTokenNotFound
	}
	if record.Status == "consumed" {
		s.logger.Info("Service", "Token already used", "token", token)
		return domain.ErrUsedToken
	}
	if time.Now().After(record.ExpiresAt) {
		if err := s.repo.UpdateUserVerificationTokenStatus(ctx, record.ID, "expired"); err != nil {
			return domain.ErrRepositoryInternalError
		}
		return domain.ErrTokenExpired
	}
	if record.Status != "pending" {
		return domain.ErrInvalidTokenState
	}

	// State verification
	err = s.repo.ConfirmVerification(ctx, record.UserID, record.ID)
	if err != nil {
		s.logger.Error("Service", "Failed to confirm verification", "error", err)
		return err
	}

	return nil
}

func (s *UserSerivce) ResendEmailVerificationToken(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	// Fetch User
	userRecord, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		s.logger.Error(domain.LogService, "Error getting user", "error", err, "email", email)
		return domain.ErrDatabaseInternalError
	}

	// Fetch Latest Verification Record
	verRecord, err := s.repo.GetVerificationByUserID(ctx, userRecord.ID)
	if err != nil { // && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Rate Limit Checks
	if verRecord != nil {
		// Strict 60s cooldown
		if time.Since(verRecord.CreatedAt) < 60*time.Second {
			return domain.ErrTooManyRequests
		}

		if verRecord.Status == "consumed" {
			return domain.ErrUserAlreadyVerified
		}
	}

	// Hourly Limit Check (Max 3 per hour)
	oneHourAgo := time.Now().Add(-1 * time.Hour)
	count, err := s.repo.GetCountsOfVerificationRecordsByUserID(ctx, userRecord.ID, oneHourAgo)
	if err != nil {
		return err
	}
	if count >= 3 {
		return domain.ErrTooManyRequests
	}

	// Prepare New Token
	token, err := GenerateSecureToken()
	if err != nil {
		s.logger.Error(domain.LogService, "Error while generating new verification token", "error", err, "user_id", verRecord.UserID)
		return domain.ErrDomainInternalError
	}

	newVer := userverification.UserVerification{
		UserID:    userRecord.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		Status:    "pending",
	}

	// Transactional Update
	// Pass verRecord.ID if it exists to invalidate it; otherwise just create
	var oldID *uuid.UUID
	if verRecord != nil {
		oldID = &verRecord.ID
	}

	return s.repo.RotateVerificationToken(ctx, *oldID, "invalidated", &newVer)
}
