package domain

import "errors"

var (
	// Domain
	ErrRepositoryInternalError = errors.New("Error while executing GetUserByEmail")
	ErrUserByEmailExists       = errors.New("User by such email is already exists")
	ErrHashingError            = errors.New("Error while hashing a password")
	ErrInvalidEmail            = errors.New("Invalid Email")
	ErrDomainInternalError     = errors.New("Domain Internal error")

	// Repository
	ErrUserAlreadyExists     = errors.New("There is no user with such email")
	ErrDatabaseInternalError = errors.New("Internal Database Error")
	ErrNotFound              = errors.New("Record not found")
	ErrTokenCollision        = errors.New("Token collision")

	// User Email Verification
	ErrTokenNotFound       = errors.New("User verification token is not found, invalid token")
	ErrTokenExpired        = errors.New("Token is expired")
	ErrInvalidTokenState   = errors.New("Invalid token")
	ErrBrokerInternalError = errors.New("Borker Internal Error")
)
