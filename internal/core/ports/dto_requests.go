package ports

type UserAndCredentialsRequest struct {
	Email        string
	UserStatus   string
	IsMFAEnabled bool
	PasswordHash string
}
