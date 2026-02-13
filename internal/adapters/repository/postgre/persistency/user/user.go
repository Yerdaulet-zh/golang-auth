package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(255);unique;index;not null"`
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
