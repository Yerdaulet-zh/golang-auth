package domainuserverification

import (
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	"github.com/google/uuid"
)

type UserVerification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	Status    string
	ExpiresAt time.Time
	CreatedAt time.Time

	User repouser.User
}
