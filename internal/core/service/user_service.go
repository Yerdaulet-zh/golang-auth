package service

import (
	"context"
	"errors"

	"github.com/golang-auth/internal/adapters/repository"
	"github.com/golang-auth/internal/core/domain"
	"github.com/golang-auth/internal/core/ports"
)

type UserSerivce struct {
	repo   *repository.UserRepository
	logger ports.Logger
}

func NewUserService(repo *repository.UserRepository, logger ports.Logger) *UserSerivce {
	return &UserSerivce{
		repo:   repo,
		logger: logger,
	}
}
func (s *UserSerivce) Register(ctx context.Context, email, password string) error {
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

// func (s *UserSerivce) CreateUser(ctx context.Context, req http.CreateUserRequest) error {
// 	_, err := s.repo.GetUserByEmail(ctx, req.Email)
// 	if err != nil {
// 		if !errors.Is(err, repository.ErrUserByEmailNotFound) {
// 			return ErrRepositoryInternalError
// 		}
// 	} else {
// 		return ErrUserByEmailExists
// 	}

// 	user_db_request := ports.UserRequest{
// 		Email: req.Email,
// 	}

// 	user_id, err := s.repo.InsertIntoUser(ctx, user_db_request)
// 	if err != nil {
// 		return repository.ErrDatabaseInternalError
// 	}

// 	user_creds_db_request := ports.UserCredentialsRequest{
// 		UserID: user_id,
// 		PasswordHash: ,
// 	}

// 	_ = user_db_request
// 	return nil

// }
