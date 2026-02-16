package userverification

import (
	"time"

	repouser "github.com/golang-auth/internal/adapters/repository/postgre/persistency/user"
	domainuserverification "github.com/golang-auth/internal/core/domain/user_verification"
	"github.com/google/uuid"
)

type UserVerification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status    string    `gorm:"type:user_verification_status;default:pending;not null"`
	ExpiresAt time.Time `gorm:"type:timestamptz;not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;default:now();not null"`

	User repouser.User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (s *UserVerification) MapToDomain() *domainuserverification.UserVerification {
	return &domainuserverification.UserVerification{
		ID:        s.ID,
		UserID:    s.UserID,
		Token:     s.Token,
		Status:    s.Status,
		ExpiresAt: s.ExpiresAt,
		CreatedAt: s.CreatedAt,
	}
}
