package usersessions

import (
	"time"

	persistency "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	"github.com/google/uuid"
)

type UserSessions struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	IPAddress  string
	UserAgent  string
	Device     string
	CreatedAt  time.Time
	LastActive time.Time
	ExpiresAt  time.Time

	User persistency.User
}
