package service

import (
	"context"
	"errors"

	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
)

type UserSerivce struct {
	repo   ports.UserRepoPorts
	logger ports.Logger
}

func NewUserService(repo ports.UserRepoPorts, logger ports.Logger) *UserSerivce {
	return &UserSerivce{
		repo:   repo,
		logger: logger,
	}
}
func (s *UserSerivce) Register(ctx context.Context, email, password string) error {
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
	repoUserCreationRequest := ports.UserAndCredentialsRequest{
		Email:        email,
		PasswordHash: hashedPassword,
	}
	userRecord, err := s.repo.GetUserByEmail(ctx, email)
	if userRecord != nil {
		return domain.ErrUserAlreadyExists
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return domain.ErrDatabaseInternalError
	}

	if err := s.repo.CreateUserWithCredentials(ctx, repoUserCreationRequest); err != nil {
		return err
	}
	return nil
}
