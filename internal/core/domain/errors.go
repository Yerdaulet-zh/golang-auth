package domain

import "errors"

var (
	// General
	ErrTooManyRequests = errors.New("Too many requests")

	// Domain
	ErrRepositoryInternalError = errors.New("Error while executing GetUserByEmail")
	ErrUserByEmailExists       = errors.New("User by such email is already exists")
	ErrHashingError            = errors.New("Error while hashing a password")
	ErrInvalidEmail            = errors.New("Invalid Email")
	ErrDomainInternalError     = errors.New("Domain Internal error")
	ErrUsedToken               = errors.New("The token already consumed")
	ErrInvaidPassword          = errors.New("Invalid password")
	ErrUserNotVerified         = errors.New("User accoutn is not verified")
	ErrUserAccountBanned       = errors.New("User account banned")
	ErrUserAccountSuspended    = errors.New("User account suspended")
	ErrTooManyUserSessions     = errors.New("Too many user sessions")
	ErrSessionNotFound         = errors.New("Session not found")

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
	ErrUserNotFound        = errors.New("User is not registered")
	ErrUserAlreadyVerified = errors.New("Email already verified/consumed")
)
