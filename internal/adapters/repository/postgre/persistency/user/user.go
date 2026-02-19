package repouser

import (
	"time"

	domain_user "github.com/golang-auth/internal/core/domain/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex:idx_email_active,where:deleted_at IS NULL;not null"`
	UserStatus   string    `gorm:"type:user_status;default:pending_verification;not null"`
	IsMFAEnabled bool      `gorm:"type:boolean;default:false;not null"`

	CreatedAt time.Time      `gorm:"type:timestamptz;default:now();not null"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;default:now();not null"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;index"`
}

type UserCredentials struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`

	LastPasswordChange time.Time `gorm:"column:last_password_change_at;type:timestamptz;default:now();not null"`
	UpdatedAt          time.Time `gorm:"type:timestamptz;default:now();not null"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func MapToDomain(u User, c UserCredentials) domain_user.User {
	return domain_user.User{
		ID:           u.ID,
		Email:        u.Email,
		UserStatus:   u.UserStatus,
		IsMFAEnabled: u.IsMFAEnabled,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt.Time, // Zero value if NULL
		Credentials: domain_user.UserCredentials{
			PasswordHash:       c.PasswordHash,
			LastPasswordChange: c.LastPasswordChange,
			UpdatedAt:          c.UpdatedAt,
		},
	}
}

func MapToORM(d domain_user.User) User {
	var gDeletedAt gorm.DeletedAt
	if !d.DeletedAt.IsZero() {
		gDeletedAt = gorm.DeletedAt{Time: d.DeletedAt, Valid: true}
	}

	return User{
		ID:           d.ID,
		Email:        d.Email,
		UserStatus:   d.UserStatus,
		IsMFAEnabled: d.IsMFAEnabled,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		DeletedAt:    gDeletedAt,
	}
}
