package ports

import "time"

type UserAndCredentialsRequest struct {
	Email                  string
	UserStatus             string
	IsMFAEnabled           bool
	PasswordHash           string
	EmailVerificationToken string
	TokenExpiration        time.Time
}
