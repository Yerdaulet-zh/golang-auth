package domain

import "errors"

var (
	// Domain
	ErrRepositoryInternalError = errors.New("Error while executing GetUserByEmail")
	ErrUserByEmailExists       = errors.New("User by such email is already exists")
	ErrHashingError            = errors.New("Error while hashing a password")

	// Repository
	ErrUserAlreadyExists     = errors.New("There is no user with such email")
	ErrDatabaseInternalError = errors.New("Internal Database Error")
	ErrNotFound              = errors.New("Record not found")
)
